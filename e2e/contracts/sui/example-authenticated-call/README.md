# SUI Authenticated WithdrawAndCall with PTB Transactions


The authenticated call works very similar to the [arbitrary call](../example-arbitrary-call/README.md), but with slight difference in PTB structure.


## PTB Structure

Instead of using three commands, the PTB for a authenticated `withdrawAndCall` transaction adds two more commands:

1. `set_message_context`:
    This command sets `sender` and `target` information right before calling the `on_call` function, so that the connected module can perform checks.

2. `reset_message_context`:
    This command clears `sender` and `target` information right after calling the `on_call` function, so that each `withdrawAndCall` uses independent message context.


```
PTB {
    // Command 0: Withdraw Implementation
    MoveCall {
        package: gateway_package_id,
        module: gateway_module,
        function: withdraw_impl,
        arguments: [
            gateway_object_ref,
            withdraw_cap_object_ref,
            coin_type,
            amount,
            nonce,
            gas_budget
        ]
    }
    
    // Command 1: Gas Budget Transfer
    TransferObjects {
        from: withdraw_impl_result[1], // Gas budget coins
        to: tss_address
    }

    // Command 2: Set Message Context
    MoveCall {
        package: gateway_package_id,
        module: gateway_module,
        function: set_message_context,
        arguments: [
            message_context,
            zevm_sender,
            target_package_id
        ]
    }

    // Command 3: Connected Module Call
    MoveCall {
        package: target_package_id,
        module: connected_module,
        function: on_call,
        arguments: [
            withdraw_impl_result[0], // Main withdrawn coins
            on_call_payload
        ]
    }

    // Command 4: Reset Message Context
    MoveCall {
        package: gateway_package_id,
        module: gateway_module,
        function: reset_message_context,
        arguments: [message_context]
    }
}
```
