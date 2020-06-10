package mutex

/*
TODO - description
 */
func receiveMutexMessage(mutexMessage MessageMutexDTO) MessageMutexDTO {
	incrementClock(mutexMessage.Time)
	// TODO
	return mutexMessage
}

/*
TODO - description
*/
func requestCriticalArea() {
	// TODO
}

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
 */
func requestCriticalArea() {
	// TODO
	// send requests
}

/*
increase local lock
 */
func incrementClock(i int32) int32 {
	clock = max(clock, i)
	return clock
}

/*
simple max function with int32 types
 */
func max(i int32, j int32) int32 {
	if i > j {
		return i
	} else {
		return j
	}
}
