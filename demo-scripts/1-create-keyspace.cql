CREATE KEYSPACE Pets_Clinic 
WITH replication = {
  'class': 'NetworkTopologyStrategy',
  'replication_factor' : 3
};
