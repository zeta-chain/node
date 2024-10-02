# Interface for arbitrary call

## **Context**

ZetaChain will support two types of smart contract calls on connected chains.

**1. Regular Calls (Authenticated Calls)**

This type of call follows the same mechanism as calls on ZetaChain. Developers need to implement a specific interface (`onCall`) in the smart contract deployed on connected chains. The `onCall` function receives a `MessageContext` that contains the sender’s address from ZetaChain.

•	This model is not yet available on mainnet; it will be supported with the next protocol version.

•	Regular calls require developers to deploy smart contracts on the connected chains, enabling the development of more advanced cross-chain applications.

•	Even though this setup requires more upfront work, the integration will be feature-complete from day one.

**2. Arbitrary Calls**

Arbitrary calls allow developers to invoke any function on any contract across connected chains without deploying new contracts on those chains.

•	This model is designed to support the chain abstraction vision, where developers don’t need to deploy contracts on connected chains.

•	However, the context of the sender from ZetaChain is lost during the smart contract call (since the gateway initiates the function call).

•	While arbitrary calls simplify cross-chain development, the feature set will be limited initially. Additional support, such as enabling NFT transfers, will be added in future updates.

## Problem

The interface is currently as follow for both model:

**Regular call**

```solidity
call(receiver, zrc20, argumentForOnCall, callOptions{}, revertOptions{})
```

**Arbitrary call**

```solidity
call(receiver, zrc20, abiEncodedFunctionCall, callOptions{isArbitraryCall: true}, revertOptions{})
```

Both models reuse the same function. An option, `isArbitraryCall`, determines whether the call is an arbitrary call or a regular call.

The third argument, bytes message, can represent different things depending on the type of call:

•	**Regular Call**: The message contains an encoded list of arguments.

•	**Arbitrary Call**: The message contains the encoded method name along with a list of arguments.

This dual-purpose usage of the same argument might be confusing for developers, as the meaning of message changes depending on the type of call.

## Alternatives to consider

### **1. Keep the Current Approach**

Maintain the existing interface, but consider renaming the argument to payload or data. These more generic terms indicate that the argument is a container for the data related to the action being performed, regardless of the call type (regular or arbitrary). This renaming can help developers better understand that the contents of this argument may vary depending on the context of the call.

**Example**

```solidity
// regular
call(
	receiver,
	zrc20,
	abi.encode(address(token), amount),
	callOptions{},
	revertOptions{},
)

// arbitrary
call(
	receiver,
	zrc20,
	abi.encodeWithSignature("swap(address,uint256)", address(token), amount),
	callOptions{isArbitraryCall: true},
	revertOptions{},
)
```

**Cons**

Even with the renaming, the argument’s dual purpose might still cause confusion for developers, as its content varies depending on whether the call is regular or arbitrary. This could lead to misunderstandings when working with the interface.

### **2. New Functions for Arbitrary Calls**

Introduce new functions specifically for arbitrary calls. These functions would include an additional argument to specify the method to be called, eliminating the need for the isArbitraryCall option.

```solidity
// regular
call(receiver, zrc20, argumentForOnCall, callOptions{}, revertOptions{})

// arbitrary
call(receiver, zrc20, methodName, arguments, callOptions{}, revertOptions{})
```

**Example**

```solidity
// arbitrary
call(
	receiver,
	zrc20,
	"swap"
	abi.encode(address(token), amount),
	callOptions{},
	revertOptions{},
)
```

**Cons**

(Stefan) i would add that user still need to pack data, but just arguments, and we would need to pack whole selector + already packed args in smart contract, which increases complexity a bit on our end, but keeps similar complexity for user

### **3 - add method name in CallOptions**

Rather than introducing a new function, the function name for arbitrary calls can be specified within the CallOptions. This change would replace the isArbitraryCall flag with a more explicit option, arbitraryCallFunctionName.

```solidity
call(receiver, zrc20, argumentForOnCall, callOptions{arbitraryCallFunctiondName: "swap"}, revertOptions{})
```

**Example**

```solidity
// arbitrary
call(
	receiver,
	zrc20,
	abi.encode(address(token), amount),
	callOptions{arbitraryCallFunctiondName: "swap"},
	revertOptions{},
)
```

**Cons**

This approach might be confusing because the method name is not just an option but a fundamental part of the call for arbitrary calls. Treating it as an option in CallOptions could obscure its importance and make it harder for developers to understand its role in the process.