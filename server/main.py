from typing import Union
from fastapi import FastAPI, Request
from milvus import collection, sentence_transformer_ef, insert_data, search_vectors
import numpy as np

app = FastAPI()



@app.get("/")
async def read_root():
    return {"Hello": "World"}


@app.post("/vector")
async def post_vector(req : Request):
    data = await req.json()
    q = data["description"]
    id = data["id"]
    if q == None or id == None:
        return {"error": "Data is not valid"}
    if len(q) == 0 or len(id) == 0 or len(q) != len(id):
        return {"error": "Data is not valid"}
    d = sentence_transformer_ef.encode_documents(q)
    embeddings = np.array(d)
    r = insert_data(collection, id, embeddings)
    print(r)
    return {"message": {"status": "success", "inserted": r.insert_count, "failed": r.err_count}}

@app.post("/query")
async def post_vector(req : Request):
    data = await req.json()
    q = data["query"]
    if q == None or len(q) == 0:
        return {"error": "Data is not valid"}
    d = sentence_transformer_ef.encode_queries(q)
    embeddings = np.array(d)
    res = search_vectors(collection, embeddings, 10)
    sd = []
    for result in res[0]:
        sd.append(result.id)
    # return {"data": sd}
    return {"message": {"status": "success", "Ids": sd}}

    