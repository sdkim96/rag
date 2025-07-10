import uvicorn

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

class ControlPlane(AgentExecutor):
    """ Control Plane for managing agents and tasks
     
    """


    @override
    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        task = context.current_task
        print("TaskID: ", task.id)
        shell = context._params.metadata.get('task_shell', None)
        print(context)
        
        
        
        
        

    @override
    async def cancel(
        self, context: RequestContext, event_queue: EventQueue
    ) -> None:
        raise Exception('cancel not supported')


ControlPlaneCard = AgentCard(
    name = 'Search Agent',
    description = 'An agent that performs search operations',
    url = 'http://localhost:8006',
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
    agent_executor=ControlPlane(),
    task_store=InMemoryTaskStore(),
)

ControlPlaneApp = A2AStarletteApplication(
    agent_card=ControlPlaneCard,
    http_handler=RequestHandler,
)


def main():
    uvicorn.run(ControlPlaneApp.build(), host="localhost", port=8006)
    

if __name__ == "__main__":
    main()

