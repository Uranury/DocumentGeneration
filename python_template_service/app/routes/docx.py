from fastapi import APIRouter, UploadFile, File, Form
from fastapi.responses import FileResponse
from app.services.docx_service import render_docx_file

router = APIRouter()

@router.post("/docx/render")
async def render_docx(template: UploadFile = File(...), data: str = Form(...)):
    import json
    data_dict = json.loads(data)  
    output_path = render_docx_file(template.file, data_dict)
    return FileResponse(
        output_path,
        media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        filename="output.docx"
    )