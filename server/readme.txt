Simple Server endpoint for the SSH-Forwarder.

As of now, since the streaming enginer is not made, im using
a mini-one for this project only to publish commands


In current version, they are batched, that is, the client
does not send the diffs in the user's terminal, rather
sends the whole command after the enter is pressed, similar
whole output is streamed back as well.
Im making using the pty backend, so that such changes can also
be easily implemented in the future when ive made the streaming
part
