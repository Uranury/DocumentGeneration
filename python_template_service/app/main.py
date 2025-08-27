from fastapi import FastAPI
from app.routes import docx, xlsx

app = FastAPI(title="Template Renderer")

app.include_router(docx.router, prefix="/docx")
app.include_router(xlsx.router, prefix="/xlsx")