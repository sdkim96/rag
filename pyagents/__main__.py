import uvicorn
from pyagents.search.main import SearchAgentServer


def main():
    
    search_agent = SearchAgentServer
    uvicorn.run(search_agent.build(), host="localhost", port=8001)

if __name__ == "__main__":
    main()