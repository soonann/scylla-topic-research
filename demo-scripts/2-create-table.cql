USE Pets_Clinic;

CREATE TABLE heartrate_v1 (
   pet_chip_id uuid,
   time timestamp,
   heart_rate int,
   PRIMARY KEY (pet_chip_id)
);
CREATE TABLE heartrate_v2 (
   pet_chip_id uuid,
   time timestamp,
   heart_rate int,
   PRIMARY KEY (pet_chip_id, time)
);
