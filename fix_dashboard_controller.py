import re

with open("internal/app/dashboard/controller/dashboard_controller.go", "r") as f:
    content = f.read()

type_map = {
    "GetOrganizerDashboard": "response.Response{data=dto.OrganizerDashboardRes}",
    "GetParticipantDashboard": "response.Response{data=dto.ParticipantDashboardRes}",
}

for func_name, success_type in type_map.items():
    pattern = r"(// @Router.*?)\nfunc \(c \*dashboardController\) " + func_name + r"\("
    
    def repl(match):
        router_line = match.group(1)
        failures = "\n// @Failure 400 {object} response.ErrorResponse\n// @Failure 401 {object} response.ErrorResponse\n// @Failure 500 {object} response.ErrorResponse\n"
        return router_line + failures + "func (c *dashboardController) " + func_name + "("
        
    content = re.sub(pattern, repl, content)

    block_pattern = re.compile(r"(// @Summary.*?\n)(.*?)(func \(c \*dashboardController\) " + func_name + r"\()", re.DOTALL)
    def repl2(match):
        comments = match.group(2)
        comments = re.sub(r"// @Success (200|201) {object} map\[string\]interface\{\}", r"// @Success \1 {object} " + success_type, comments)
        return match.group(1) + comments + match.group(3)
    content = block_pattern.sub(repl2, content)

with open("internal/app/dashboard/controller/dashboard_controller.go", "w") as f:
    f.write(content)
