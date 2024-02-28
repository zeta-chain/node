# Messages

## MsgWithdrawEmission

WithdrawEmission create a withdraw emission object , which is then process at endblock
The withdraw emission object is created and stored
using the address of the creator as the index key ,therefore, if more that one withdraw requests are created in a block on thr last one would be processed.
Creating a withdraw does not guarantee that the emission will be processed
All withdraws for a block are deleted at the end of the block irrespective of whether they were processed or not.

```proto
message MsgWithdrawEmission {
	string creator = 1;
	string amount = 2;
}
```

