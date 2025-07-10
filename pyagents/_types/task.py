from typing import List
from pydantic import  BaseModel, Field

from a2a.types import Task

class Trip(BaseModel):
    departure: str = Field(..., description="The departure location of the task.")
    destination: str = Field(..., description="The destination location of the task.")
    trip_count: int = Field(..., description="The number of trips made by the task.")

class TaskShell(BaseModel):
    task_id: str = Field(..., description="The task details.")
    context_id: str = Field(..., description="The context ID associated with the task.")
    trip: Trip = Field(..., description="The trip details associated with the task.")
    
