USE Pets_Clinic;

SELECT * 
FROM heartrate_v2 
WHERE 
  pet_chip_id = 123e4567-e89b-12d3-a456-426655440b23 
  AND time >='2019-03-04 07:01:00';

SELECT * 
FROM heartrate_v1 
WHERE 
  pet_chip_id = 123e4567-e89b-12d3-a456-426655440b23 
  AND time >='2019-03-04 07:01:00';
