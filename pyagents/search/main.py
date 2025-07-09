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


from typing_extensions import override

async def some_information_from_vectorstore():
    yield 'Seoul is the capital of South Korea.'
    await asyncio.sleep(2)
    yield 'About 10million people live in Seoul.'

class SearchAgent(AgentExecutor):
    """ Search for vectorstore """


    @override
    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        
        task = context.current_task
        if not task:
            task = new_task(context.message) # type: ignore
            await event_queue.enqueue_event(task)

        updater = TaskUpdater(event_queue, task.id, task.contextId)
        query = context.get_user_input()
        
        await updater.update_status(
            state=TaskState.working,
            message=new_agent_text_message(
                "🔍 Searching for valuable information..",
                task.contextId,
                task.id,
            )
        )
        await asyncio.sleep(0.1)
        
        async for part in some_information_from_vectorstore():
            await updater.add_artifact(
                [Part(root=TextPart(text=part))],
                name='conversion_result',
            )

        await updater.update_status(
            state=TaskState.completed,
            message=new_agent_text_message(
                "🔍 Thanks for your visit!",
                task.contextId,
                task.id,
            ),
            final=True
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
