import uvicorn

from a2a.server.apps import A2AStarletteApplication
from a2a.server.tasks import InMemoryTaskStore
from a2a.server.request_handlers import DefaultRequestHandler

from pyagents.search.main import SearchAgentServer
from pyagents.entrypoint import EntrypointCard, Entrypoint
from pyagents.control_plane import ControlPlane, ControlPlaneCard


EntrypointHandler = DefaultRequestHandler(
    agent_executor=Entrypoint(control_plane=ControlPlaneCard),
    task_store=InMemoryTaskStore(),
)

EntrypointApp = A2AStarletteApplication(
    agent_card=EntrypointCard,
    http_handler=EntrypointHandler,
)


def main():
    
    search_agent = SearchAgentServer
    uvicorn.run(search_agent.build(), host="localhost", port=8001)

if __name__ == "__main__":
    main()