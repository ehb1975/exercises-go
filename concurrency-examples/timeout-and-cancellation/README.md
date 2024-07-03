# Timeout and Cancellation

- Timeout and cancellation: Channels can be used to implement timeouts and cancellation of long-running operations. By using a select statement with a time. After the channel, you can wait for a specified amount of time before proceeding with an operation. Additionally, by using a cancel channel, you can signal to a goroutine that it should stop processing and return early.

Source: https://medium.com/@varmapooja09/mastering-go-channels-part-1-7baa978a7de8

