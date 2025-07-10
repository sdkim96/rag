import uuid
import httpx

from a2a.server.apps import A2AStarletteApplication
from a2a.server.tasks import InMemoryTaskStore
from a2a.server.request_handlers import DefaultRequestHandler
from a2a.types import (
    AgentCapabilities,
    AgentCard,
    AgentSkill,
    AgentCard,
    MessageSendParams,
    SendMessageRequest,
    SendStreamingMessageRequest,
)
from a2a.server.agent_execution import AgentExecutor, RequestContext
from a2a.server.events import EventQueue
from a2a.utils import (
    new_agent_text_message,
    new_task,
)
from a2a.client import A2ACardResolver, A2AClient

from pyagents._types.task import TaskShell, Trip
from typing_extensions import override

class Entrypoint(AgentExecutor):
    """ Control Plane for managing agents and tasks
     
    """

    def __init__(
        self,
        control_plane: AgentCard
    ):
        self.control_plane = control_plane


    @override
    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        task = context.current_task
        message = new_agent_text_message(
            context.get_user_input(),
            context_id=context.context_id,
            task_id=context.task_id,
        )
        if task is None:
            task = new_task(message)

        shell = TaskShell(
            task_id=task.id,
            trip=Trip(
                departure='entrypoint',
                destination='control_plane',
                trip_count=1
            )
        )
        async with httpx.AsyncClient(timeout=70) as httpx_client:
            client = A2AClient(
                httpx_client=httpx_client, agent_card=self.control_plane
            )
            request = SendMessageRequest(
                id=str(uuid.uuid4()), 
                params=MessageSendParams(
                    message=message,
                    metadata={
                        "shell": shell.model_dump(mode='json', exclude_none=True),
                    }
                )
            )

            response = await client.send_message(request)
            
        

            

    @override
    async def cancel(
        self, context: RequestContext, event_queue: EventQueue
    ) -> None:
        raise Exception('cancel not supported')


EntrypointCard = AgentCard(
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
