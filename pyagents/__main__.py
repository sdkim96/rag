import uvicorn

from a2a.server.apps import A2AStarletteApplication
from a2a.server.request_handlers import DefaultRequestHandler
from a2a.server.tasks import InMemoryTaskStore
from a2a.types import (
    AgentCapabilities,
    AgentCard,
    AgentSkill,
)

from pyagents.search.main import SearchAgentServer


def main():
    
    search_agent = SearchAgentServer

    uvicorn.run(search_agent.build(), host="localhost", port=8001)

if __name__ == "__main__":
    main()