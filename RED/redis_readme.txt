



// LIST FUNCTIONS
//  _       _________ _______ _________ _______
// ( \      \__   __/(  ____ \\__   __/(  ____ \
// | (         ) (   | (    \/   ) (   | (    \/
// | |         | |   | (_____    | |   | (_____
// | |         | |   (_____  )   | |   (_____  )
// | |         | |         ) |   | |         ) |
// | (____/\___) (___/\____) |   | |   /\____) |
// (_______/\_______/\_______)   )_(   \_______)

// Txt2List, List_RandItem

// convert a textfile into a redis list

//  _   _  _____ _____ _____
// | | | |/  ___|  ___|_   _|
// | |_| |\ `--.| |__   | |
// |  _  | `--. \  __|  | |
// | | | |/\__/ / |___  | |
// \_| |_/\____/\____/  \_/
//
// Store object-like data : { key1:1value, key2:value2...}
// 🟥 HGETALL h_test             -> return all keys and values
// 🟥 HVALS h_test               -> return all vals
// 🟥 HGET h_test c1             -> return h_test[c1] (value)
// 🟥 HSET h_test field1 "Hello" -> add {field1:"Hello"}
// 🟥 HINCRBY h_test joon 1      -> add 1 HSET (one str at a time)
// 🟥 HMGET h_test field1 field2 -> get multiple field values

//  ______ _____ _____ _____
// |___  //  ___|  ___|_   _|
//    / / \ `--.| |__   | |
//   / /   `--. \  __|  | |
// ./ /___/\__/ / |___  | |
// \_____/\____/\____/  \_/
//
// ZSET --Store unique, ordered collections with scores (rankings)
// 🟥 ZRANGE key start stop [WITHSCORES]: Get members by index (rank).
// 🟥 ZRANGEBYSCORE key min max [WITHSCORES]: Get members by score range.
// 🟥 ZREM key member: Remove a member.
// 🟥 ZSCORE key member: Get the score of a specific member.
// 🟥 ZRANK key member: Get the rank (0-based index) of a member.
// 🟥 ZINCRBY key increment member: Increment the score of a member