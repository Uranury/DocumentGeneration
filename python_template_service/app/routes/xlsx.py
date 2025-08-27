from fastapi import APIRouter, UploadFile, File
from fastapi.responses import FileResponse
from app.services.xlsx_service import render_xlsx_file

router = APIRouter()

@router.post("/render")
async def render_xlsx(template: UploadFile = File(...), data: str = File(...)):
    output_path = render_xlsx_file(template.file, data)
    return FileResponse(
        output_path,
        media_type="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        filename="output.xlsx"
    )
