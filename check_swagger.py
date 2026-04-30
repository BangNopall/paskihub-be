import glob
import re

files = glob.glob('internal/app/**/*_controller.go', recursive=True)

missing_swagger = []

for file in files:
    with open(file, 'r') as f:
        lines = f.readlines()
    
    for i, line in enumerate(lines):
        if line.startswith('func ('):
            # Try to extract the method name. It expects: func (receiver) MethodName(...)
            match = re.match(r'func \([^)]+\)\s+([A-Z]\w*)', line)
            if match:
                func_name = match.group(1)
                
                if func_name in ['Route', 'Init', 'Setup']:
                    continue
                
                has_swagger = False
                for j in range(i-1, -1, -1):
                    prev_line = lines[j].strip()
                    if not prev_line.startswith('//'):
                        break
                    if '@Summary' in prev_line or '@Router' in prev_line:
                        has_swagger = True
                        break
                
                if not has_swagger:
                    missing_swagger.append(f"{file}:{i+1} - {func_name}")

if missing_swagger:
    print("Endpoints missing Swagger:")
    for m in missing_swagger:
        print(m)
else:
    print("All exported controller methods have Swagger.")
