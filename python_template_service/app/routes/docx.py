from fastapi import APIRouter, UploadFile, File
from fastapi.responses import FileResponse
from app.services.docx_service import render_docx_file

router = APIRouter()

@router.post("/render")
async def render_docx(template: UploadFile = File(...), data: str = File(...)):
    output_path = render_docx_file(template.file, data)
    return FileResponse(output_path, media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document", filename="output.docx")
