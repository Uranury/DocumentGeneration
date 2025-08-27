import tempfile
import shutil
import json
from xlsx_template import Template

def render_xlsx_file(template_file, json_data_str):
    json_data = json.loads(json_data_str)

    # Save uploaded template to a temporary file
    with tempfile.NamedTemporaryFile(delete=False, suffix=".xlsx") as tmp_template:
        shutil.copyfileobj(template_file, tmp_template)
        tmp_template_path = tmp_template.name

    tpl = Template(tmp_template_path)
    tpl.render(json_data)

    # Save rendered XLSX to temp file
    output_file = tempfile.NamedTemporaryFile(delete=False, suffix=".xlsx")
    tpl.save(output_file.name)
    return output_file.name
