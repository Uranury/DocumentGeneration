# FastAPI Router
from fastapi import APIRouter, UploadFile, File, Form
from fastapi.responses import FileResponse
from app.services.xlsx_service import render_xlsx_file

router = APIRouter()

@router.post("/xlsx/render")
async def render_xlsx(template: UploadFile = File(...), data: str = Form(...)):
    import json
    data_dict = json.loads(data)  
    output_path = render_xlsx_file(template.file, data_dict)
    return FileResponse(
        output_path,
        media_type="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        filename="output.xlsx"
    )
