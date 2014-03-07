package utils

func gen_id(timestamp int64, machineid int64, sequence_number int64) int64{

  return timestamp << 22 | machineid << 12 | sequence_number

}
