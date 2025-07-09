import asyncio
import uvicorn

from a2a.server.apps import A2AStarletteApplication
from a2a.client import A2AClient
from a2a.server.tasks import InMemoryTaskStore
from a2a.server.request_handlers import DefaultRequestHandler
from a2a.types import (
    AgentCapabilities,
    AgentCard,
    AgentSkill,
    Message,
    Part,
    TextPart,
    Role,
    SendMessageRequest
)
from a2a.server.agent_execution import AgentExecutor, RequestContext
from a2a.server.events import EventQueue
from a2a.utils import new_agent_text_message


from typing_extensions import override

async def some_gen():
    yield 'hello'
    await asyncio.sleep(1)
    yield 'world'
    
class HelloWorldAgentExecutor(AgentExecutor):
    """Test AgentProxy Implementation."""


    @override
    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        
        async for part in some_gen():
            await event_queue.enqueue_event(new_agent_text_message(part))
        

    @override
    async def cancel(
        self, context: RequestContext, event_queue: EventQueue
    ) -> None:
        raise Exception('cancel not supported')


SearchAgentCard = AgentCard(
    name = 'Search Agent',
    description = 'An agent that performs search operations',
    url = 'http://localhost:8002',
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

request_handler = DefaultRequestHandler(
    agent_executor=HelloWorldAgentExecutor(),
    task_store=InMemoryTaskStore(),
)

SearchAgentServer = A2AStarletteApplication(
    agent_card=SearchAgentCard,
    http_handler=request_handler,
)
