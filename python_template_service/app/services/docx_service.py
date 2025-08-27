import tempfile
import shutil
import json
from docxtpl import DocxTemplate

def render_docx_file(template_file, data_dict):
    from docxtpl import DocxTemplate
    import tempfile, shutil

    with tempfile.NamedTemporaryFile(delete=False, suffix=".docx") as tmp_template:
        shutil.copyfileobj(template_file, tmp_template)
        tmp_template_path = tmp_template.name

    doc = DocxTemplate(tmp_template_path)
    doc.render(data_dict)  # already a dict

    output_file = tempfile.NamedTemporaryFile(delete=False, suffix=".docx")
    doc.save(output_file.name)
    return output_file.name
