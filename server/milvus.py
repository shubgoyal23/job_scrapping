from pymilvus import connections, FieldSchema, CollectionSchema, DataType, Collection
from dotenv import load_dotenv
import os

from pymilvus import model
load_dotenv()



sentence_transformer_ef = model.dense.SentenceTransformerEmbeddingFunction(
    model_name='multi-qa-mpnet-base-cos-v1', 
    device='cpu'
    )


host= os.getenv('MILVUSHOST')
token= os.getenv('MILVUSTOKEN')
db_name=os.getenv('MILVUS_DBNAME')

def connect_to_milvus():
    connections.connect(
        "default",
        uri=host, 
        token=token, 
        db_name=db_name, 
        )


def create_collection(collection_name, dim):
    """
    Create a Milvus collection with an integer primary key and a vector field.
    """
    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=dim)
    ]
    schema = CollectionSchema(fields, description="Collection for job vector embeddings")
    collection = Collection(name=collection_name, schema=schema)
    print(f"Collection '{collection_name}' created successfully!")
    return collection

def insert_data(collection, ids, embeddings):
    """
    Insert vector data into the collection.
    """
    data = [ids, embeddings.tolist()]
    res = collection.insert(data)
    return res
def search_vectors(collection, query_vector, top_k=3):
    """
    Search for similar vectors in the collection.
    """
    # collection.load()
    results = collection.search(
        data=query_vector.tolist(),
        anns_field="embedding",
        param={"metric_type": "COSINE", "params": {"nprobe": 10}},
        limit=top_k
    )
    print("Search results:")
    for result in results[0]:
        print(f"ID: {result.id}, Distance: {result.distance}")
    return results
        
connect_to_milvus()
collection = create_collection("jobVector", 768)