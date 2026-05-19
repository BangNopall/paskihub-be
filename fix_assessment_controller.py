import re

with open("internal/app/assessment/controller/assessment_controller.go", "r") as f:
    content = f.read()

# Map of function name to its success type
type_map = {
    "CreateJudge": "response.Response{data=dto.JudgeRes}",
    "GetJudges": "response.Response{data=[]dto.JudgeRes}",
    "UpdateJudge": "response.Response{data=dto.JudgeRes}",
    "DeleteJudge": "response.Response",
    "CreateViolationType": "response.Response{data=dto.ViolationTypeRes}",
    "GetViolationTypes": "response.Response{data=[]dto.ViolationTypeRes}",
    "UpdateViolationType": "response.Response{data=dto.ViolationTypeRes}",
    "DeleteViolationType": "response.Response",
    "CreateScoreCategory": "response.Response{data=dto.ScoreCategoryRes}",
    "GetScoreCategories": "response.Response{data=[]dto.ScoreCategoryRes}",
    "UpdateScoreCategory": "response.Response{data=dto.ScoreCategoryRes}",
    "DeleteScoreCategory": "response.Response",
    "CreateScoreSubCategory": "response.Response{data=dto.ScoreSubCategoryRes}",
    "UpdateScoreSubCategory": "response.Response{data=dto.ScoreSubCategoryRes}",
    "DeleteScoreSubCategory": "response.Response",
    "GetUnifiedAssessment": "response.Response{data=dto.UnifiedAssessmentRes}",
    "InputScore": "response.Response{data=dto.ScoreRes}",
    "CreateAward": "response.Response{data=dto.AwardRes}",
    "GetAwards": "response.Response{data=[]dto.AwardRes}",
    "UpdateAward": "response.Response{data=dto.AwardRes}",
    "DeleteAward": "response.Response",
}

for func_name, success_type in type_map.items():
    # Find the block of comments preceding the function
    pattern = r"(// @Router.*?)\nfunc \(c \*assessmentController\) " + func_name + r"\("
    
    def repl(match):
        router_line = match.group(1)
        failures = "\n// @Failure 400 {object} response.ErrorResponse\n// @Failure 401 {object} response.ErrorResponse\n// @Failure 500 {object} response.ErrorResponse\n"
        return router_line + failures + "func (c *assessmentController) " + func_name + "("
        
    content = re.sub(pattern, repl, content)

    # replace Success object map[string]interface{} with the success_type for this specific function block
    # We do a targeted replace by finding the function comment block
    block_pattern = re.compile(r"(// @Summary.*?\n)(.*?)(func \(c \*assessmentController\) " + func_name + r"\()", re.DOTALL)
    def repl2(match):
        comments = match.group(2)
        # replace Success
        comments = re.sub(r"// @Success (200|201) {object} map\[string\]interface\{\}", r"// @Success \1 {object} " + success_type, comments)
        return match.group(1) + comments + match.group(3)
    content = block_pattern.sub(repl2, content)

with open("internal/app/assessment/controller/assessment_controller.go", "w") as f:
    f.write(content)
