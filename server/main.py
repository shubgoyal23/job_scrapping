from typing import Union

from fastapi import FastAPI, Request
# from sbert import Get_embeddings

from milvus import sentence_transformer_ef, create_collection, connect_to_milvus, insert_data


app = FastAPI()
connect_to_milvus()
collection = create_collection("jobVector", 768)


@app.get("/")
async def read_root():
    return {"Hello": "World"}


@app.get("/items/{item_id}")
async def read_item(item_id: int, q: Union[str, None] = None):
    return {"item_id": item_id, "q": q}

@app.post("/vector")
async def post_vector(req : Request):
    data = await req.json()
    q = data["query"]
    id = data["id"]
    if q == None or id == None:
        return {"data": "error"}
    # d = Get_embeddings(q)
    d = sentence_transformer_ef.encode_documents(q)
    insert_data(collection, id, d)
    # print(sentence_transformer_ef.dim, d[0].shape)
    return {"data": "d"}

@app.post("/query")
async def post_vector(req : Request):
    data = await req.json()
    q = data["query"]
    # d = Get_embeddings(q)
    d = sentence_transformer_ef.encode_queries(q)
    print(sentence_transformer_ef.dim, d[0].shape)
    return {"data": "d"}
    
    