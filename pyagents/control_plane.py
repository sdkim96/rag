import asyncio

from a2a.server.apps import A2AStarletteApplication
from a2a.server.tasks import InMemoryTaskStore
from a2a.server.request_handlers import DefaultRequestHandler
from a2a.types import (
    AgentCapabilities,
    AgentCard,
    AgentSkill,
    TaskState,
    Part,
    TextPart,
)
from a2a.server.agent_execution import AgentExecutor, RequestContext
from a2a.server.events import EventQueue
from a2a.utils import (
    new_agent_text_message,
    new_task,
)
from a2a.server.tasks import TaskUpdater

from pyagents._types.task import TaskPod
from typing_extensions import override

class ControlPlane(AgentExecutor):
    """ Control Plane for managing agents and tasks
     
    """


    @override
    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        
        pod = TaskPod(
            
        )
        
        
        
        
        

    @override
    async def cancel(
        self, context: RequestContext, event_queue: EventQueue
    ) -> None:
        raise Exception('cancel not supported')


SearchAgentCard = AgentCard(
    name = 'Search Agent',
    description = 'An agent that performs search operations',
    url = 'http://localhost:8001',
    version = '1.0.0',
    defaultInputModes = ['text'],
    defaultOutputModes = ['text'],
    capabilities = AgentCapabilities(streaming=True),
    skills = [
        AgentSkill(
            id='01',
            tags=['search'],
            name='search',
            description='Performs a search operation',
            inputModes=['text'],
            outputModes=['text'],
        ),
    ],
)

RequestHandler = DefaultRequestHandler(
    agent_executor=SearchAgent(),
    task_store=InMemoryTaskStore(),
)

SearchAgentServer = A2AStarletteApplication(
    agent_card=SearchAgentCard,
    http_handler=RequestHandler,
)
