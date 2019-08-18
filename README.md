# Roadie

Roadie facilitates an atomic swap between two parties that would like to
exchange siacoins for ether. Roadie was created to help set the stage for
[Sia](https://sia.tech/).

![Roadie demo](https://i.imgur.com/Hfdl9B9.gif)
(The screencast has been sped up and edited for clarity.)

## How It Works

Alice has ether on the Ethereum blockchain and would like to trade with Bob, who
has siacoins on the Sia blockchain. Bob is running a server, which he has
registered with the roadie smart contract on Ethereum. Alice looks up the
address of Bob's server and connects to it with a request detailing how many
siacoins she would like to buy. She receives an offer and decides to accept it.

Sia uses Schnorr signatures, which have the nice property that they support
["native multisig"](https://bitcoincore.org/en/2017/03/23/schnorr-signature-aggregation/).
This means that Alice and Bob can each create a fresh public/private key pair
and combine the public keys to create a 2-of-2 multisig Sia address. (Sia also
has built-in support for multisig, but this feature is not used here.) Bob
transfers siacoins to this address, but not before creating a refund transaction
and having it signed by Alice. This refund transaction will transfer all
siacoins back to Bob, but is timelocked and can only be used in the near future.
Should anything go wrong with the atomic swap, Bob can simply wait and then
publish the refund transaction.

The atomic swap itself makes uses of an "adaptor" signature scheme (as described
in ["Flipping the scriptless script on
Schnorr"](https://joinmarket.me/blog/blog/flipping-the-scriptless-script-on-schnorr/)).
Bob prepares a transaction which moves all siacoins from the 2-of-2 address to
Alice. But instead of signing and publishing it in the regular way - which would
allow Alice to just take the coins - he instead creates an adaptor
public/private key pair and publishes an adaptor signature for the payout
transaction. Alice cannot use this signature directly, but she can verify that
if she were to learn the adaptor secret, she would be able to claim the coins.

Alice now locks up her ether in the roadie smart contract, where it can be
released only under one of two conditions. The first is a timelock again - if
anything goes wrong, Alice can simply wait and then reclaim her ether. The
second condition is for Bob to reveal the adaptor secret and receive the
payment.

Bob now sends the adaptor secret to the smart contract. The smart contract
checks that it is indeed the correct secret (using 
[ed25519-solidity](https://github.com/javgh/ed25519-solidity), an implemention
of the necessary elliptic curve math in Solidity to check that the adaptor
secret matches the adaptor public key) and releases the ether payment to Bob.
Alice uses the adaptor secret to build a valid payout transaction and claims the
siacoins from the 2-of-2 multisig. The atomic swap is complete.

## Installation

Roadie needs a [Sia node](https://sia.tech/get-started) with an unlocked wallet
and an Ethereum node (for example [geth](https://github.com/ethereum/go-ethereum)).
Your distribution might have packages for this software. It is all written in
Go, so you can also use `go get` to install it. In this case the software will -
by default - be installed in `~/go/bin/`, so make sure that this directory is on
your `$PATH`:

    $ export PATH=$PATH:~/go/bin

Installing and running Sia:

    $ go get -u gitlab.com/NebulousLabs/Sia/...     # '...' is important
    $ mkdir sia; cd sia
    ~/sia$ siad
    ~/sia$ siac wallet init
    ~/sia$ siac wallet unlock

Installing and running geth:

    $ go get -u github.com/ethereum/go-ethereum/cmd/geth
    $ geth --syncmode light

Installing roadie:

    $ go get -u github.com/javgh/roadie/cmd/roadie

## Usage

Roadie uses its own Ethereum wallet. After initialization, the wallet will be
stored in `~/.config/roadie/keystore` with the password set to blank. If
necessary, it can also be read by `geth` and many other Ethereum wallets.

    $ roadie init
    Creating new Ethereum keystore ~/.config/roadie/keystore
    Ethereum address: <new Ethereum address here>
    Ethereum balance: 0.000000 ETH

    Please deposit funds into the address listed above.
    A minimum of 0.010000 ETH is needed to proceed.

After funding the wallet, you are ready to buy siacoins:

    $ roadie buy 1

See `roadie help` for additional options. The command `roadie serve` is
currently for advanced users only. Not all aspects of running a server are
documented yet. It will also probably be necessary to implement a custom pricing
strategy - see `FixedPremiumTrader` in `trader/trader.go` for an example.

## Sequence Diagram

    Alice                                   Bob
    -----                                   ---
                                            register server on blockchain
    request non-binding offer for X SC
                                            send non-binding offer including:
                                            message, availability, X ether, X anti-spam fee
    decide on offer
    pick fresh anti-spam id
    make anti-spam payment to
      hash(anti-spam id) and
      wait for confirmations
    request binding offer with
      anti-spam id
                                            verify anti-spam payment
                                            lock in binding offer
                                            if lock successful: remember anti-spam id
                                            send binding offer including:
                                            message, availability, X ether
    decide on offer
    accept offer and
      send alicePubKey
                                            build funding tx
                                            build refund tx
                                              with timelock 4 hours
                                            send bobPubkey, fundingOutputID,
                                              bobRefundUnlockHash, timelock,
                                              bobRefundNoncePoint
    check timelock
    build refund tx
    sign refund tx
    send aliceRefundNoncePoint,
      refundSigAlice
                                            sign refund tx
                                            verify refund tx
                                            sign funding tx
                                            broadcast funding tx
                                            keep refund tx for potential rollback
                                            send funding tx id
    wait for funding tx id
    verify funding tx
    build claim tx
    send aliceClaimUnlockHash,
      aliceClaimNoncePoint
                                            build claim tx
                                            generate adaptor
                                            adaptor sign claim tx
                                            send bobClaimNoncePoint,
                                              adaptorPubKey, adaptorSigBob,
                                              depositRecipient
    verify adaptor sig
    deposit ether for adaptor
      with timelock 2 hours
      and burn anti-spam id
    remember details for
      potential rollback
    wait for deposit tx
    announce deposit
                                            check deposit
                                            claim deposit and publish adaptor
                                            send ok
    wait for adaptor
    adaptor sign claim tx
    combine sigs and adaptor
    broadcast claim tx


## Gas Cost

The smart contract unfortunately does require a fair amount of gas. Especially
performing the necessary elliptic curve math can use more than 1 000 000 gas in
some cases. It should be possible to optimize this further. Ethereum might add
additional precompiles for these kinds of operations or it might be possible to
port the optimized code at [http://ed25519.cr.yp.to/](http://ed25519.cr.yp.to/)
to Solidity.

## Acknowledgments

This project makes heavy use of prior work done by the Hyperspace developers on
a [manual atomic swap](https://github.com/HyperspaceApp/atomicswap).
