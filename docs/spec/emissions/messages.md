# Messages

## MsgWithdrawEmission

WithdrawEmission allows the user to withdraw from their withdrawable emissions.
on a successful withdrawal, the amount is transferred from the undistributed rewards pool to the user's account.
if the amount to be withdrawn is greater than the available withdrawable emission, the max available amount is withdrawn.
if the pool does not have enough balance to process this request, an error is returned.

```proto
message MsgWithdrawEmission {
	string creator = 1;
	string amount = 2;
}
```

