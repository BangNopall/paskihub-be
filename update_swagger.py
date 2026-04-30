import os
import glob

for filepath in glob.glob('/Users/noxval/_PROJECT_/paskihub-be/internal/app/**/*_controller.go', recursive=True):
    with open(filepath, 'r') as f:
        lines = f.readlines()
    
    new_lines = []
    in_doc_block = False
    has_security = False
    
    for line in lines:
        if ' godoc' in line or '// @Summary' in line:
            if not in_doc_block:
                in_doc_block = True
                has_security = False
                
        if in_doc_block:
            if '// @Security' in line:
                has_security = True
                if 'ApiKeyAuth' not in line:
                    if 'BearerAuth' in line:
                        line = line.replace('BearerAuth', 'ApiKeyAuth && BearerAuth')
                    else:
                        # Append it
                        line = line.replace('// @Security', '// @Security ApiKeyAuth &&')
            if '// @Router' in line:
                if not has_security:
                    # Insert before this line
                    indent = line[:line.find('// @Router')]
                    new_lines.append(indent + '// @Security ApiKeyAuth\n')
                in_doc_block = False
                has_security = False
                
        new_lines.append(line)
        
    with open(filepath, 'w') as f:
        f.writelines(new_lines)

print("Done")
