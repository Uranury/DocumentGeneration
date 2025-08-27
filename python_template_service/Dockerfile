# Use a lightweight Python base image
FROM python:3.12-slim

# Set working directory
WORKDIR /app

# Copy requirements and install them
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the app source code
COPY app/ ./app

# Expose port for FastAPI
EXPOSE 8000

# Start the service with uvicorn
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
