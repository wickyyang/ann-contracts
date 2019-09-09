pragma solidity ^0.5.0;

import "github.com/OpenZeppelin/openzeppelin-contracts/contracts/math/SafeMath.sol";

/**
 * @title Multisig
 * @dev Multisignature Sample Contract.
 */
contract Multisig {
    using SafeMath for uint256;
    
    // Contract's administrator
    address private _admin;
    // The number of owners that must confirm the same operation before it is run.
    uint256 private _m_required;
    // Pointer used to find a free slot in m_owners
    uint256 private _m_numOwners;
    // List of owners, limit is 10.
    address[10] private _m_owners;
    // Index on the list of owners to allow reverse lookup
    mapping(uint256 => uint256) private _m_ownerIndex;
    // The ongoing operations.
    mapping(bytes32 => PendingState) private _m_pending;

    // Struct for the status of a pending operation.
    struct PendingState {
        uint256 yetNeeded;
        uint256 ownersDone;
        uint256 index;
        address[10] ownersAddr;
    }

    event OwnerInit(address[] indexed newOwner);
    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event OwnerAdded(address indexed newOwner);
    event OwnerRemoved(address indexed oldOwner);
    event RequirementChanged(uint256 newRequirement);
    event Confirmation(address indexed owner, bytes32 operation);
    event ChangeAdmin(address indexed admin);

    constructor() public {
        _admin = msg.sender;
    }

    modifier onlyAdmin {
        require (msg.sender == _admin);
        _;
    }

    /**
     * @dev Initializes the multisignature administrators and 
     * the number of signatures needs to be collectted.
     * @param _owners Multisignature administrator's address.
     * @param _required The number of signatures needs to be collectted.
     */
    function multiOwner(address[] memory _owners, uint256 _required) onlyAdmin public {
        _m_numOwners = _owners.length.add(1);
        require (_m_numOwners >= _required);
        
        _m_owners[0] = msg.sender;
        _m_ownerIndex[uint256(msg.sender)] = 0;
        for (uint256 i = 0; i < _owners.length; i++) {
            _m_owners[i.add(1)] = _owners[i];
            _m_ownerIndex[uint256(_owners[i])] = i.add(1);
        }

        _m_required = _required;

       emit OwnerInit(_owners);
    }

    /**
     * @dev Changes the multisignature administrator.
     * @param _from the original multisignature administrator's address.
     * @param _to the new multisignature administrator's address.
     */ 
    function changeOwner(address _from, address _to) onlyAdmin public {
        require (_isSigner(_from) && !_isSigner(_to));

        uint256 ownerIndex = _m_ownerIndex[uint256(_from)];
        _m_owners[ownerIndex] = _to;
        delete _m_ownerIndex[uint256(_from)];
        _m_ownerIndex[uint256(_to)] = ownerIndex;

        emit OwnerChanged(_from, _to);
    }

    /**
     * @dev Adds new multi-signature administrator.
     * @param _owner the new multi-signature administrator's address.
     */ 
    function addOwner(address _owner) onlyAdmin public {
        require (!_isSigner(_owner) && _owner != address(0));

        for (uint256 i = 0; i < _m_owners.length; i++) {
            if (_m_owners[i] == address(0)) {
                _m_owners[i] = _owner;
                _m_ownerIndex[uint256(_owner)] = i;

                _m_numOwners++;
                emit OwnerAdded(_owner);
                break;
            }
        }
    }

    /**
     * @dev Deletes the multi-signature administrator
     * @param _owner the multi-signature administrator's address needs to be deleted.
     */  
    function removeOwner(address _owner) onlyAdmin public {
        require (_isSigner(_owner));

        uint256 ownerIndex = _m_ownerIndex[uint256(_owner)];
        delete _m_owners[ownerIndex];
        delete _m_ownerIndex[uint256(_owner)];
        _m_numOwners --;

        emit OwnerRemoved(_owner);
    }

    /**
     * @dev Tells whether is a multisignature administrator.
     * @param _owner the owner's address needed to be checked.
     * @return the result
     */
    function checkOwner(address _owner) public view returns(bool) {
        return _isSigner(_owner);
    }

    /**
     * @dev Changes the number of signatures needs to be collectted.
     * @param _newRequired new number of signatures needs to be collectted.
     */ 
    function changeRequirement(uint256 _newRequired) onlyAdmin public {
        require (_newRequired <= _m_numOwners);
        _m_required = _newRequired;

        emit RequirementChanged(_newRequired);
    }

    /**
     * @dev Tells whether is a multi-signature administrator.
     * @param _signer the owner's address needed to be checked.
     * @return the result
     */ 
    function _isSigner(address _signer) private view returns (bool) {
        if (_signer == address(0)) return false;
        // Iterate through all signers on the wallet and
        for (uint256 i = 0; i < _m_owners.length; i++) {
          if (address(_m_owners[i]) == _signer) {
            return true;
          }
        }
        
        return false;
    }

    /**
     * @dev Verify the signature
     * @param _contractAddr contract's address
     * @param _from transfer's address
     * @param _to recipient's address
     * @param _value the amount of tokens to be spent
     * @param _signature the multi-signature administrator's signature
     * @return the verification results. If the signature is right and 
     * the number of signatures needs to be collectted is enough, return true.
     */ 
    function confirmAndCheck(address _contractAddr, address _from, address _to, uint256 _value, bytes memory _signature) onlyAdmin public returns (bool) {
        bytes32 operationHash = keccak256(
            abi.encodePacked(_contractAddr, _from, _to, _value)
        );

        if (!_verifyMultiSig(operationHash, _signature)) {
            return false;
        }

        if (_checkSigNum(operationHash, _signature)) {
            return true;
        } 
    }

    /**
     * @dev Verify the signature
     * @param _operationHash operation hash
     * @param _signature the multi-signature administrator's signature
     * @return the verification results. If the signature is right, return true.
     */ 
    function _verifyMultiSig(bytes32 _operationHash, bytes memory _signature) private view returns (bool) {

    address otherSigner = _recoverAddressFromSignature(_operationHash, _signature);

    if (!_isSigner(otherSigner)) {
        // Other signer not on this wallet or operation does not match arguments
        return false;
    }

    if (otherSigner == msg.sender) {
        // Cannot approve own transaction
        return false;
    }

        return true;
    }

    /**
     * @dev Gets signer's address using ecrecover.
     * @param _operationHash operation hash
     * @param _signature the multi-signature administrator's signature
     * @return address recovered from the signature
     */
    function _recoverAddressFromSignature(bytes32 _operationHash, bytes memory _signature) private pure returns (address) {
        require (_signature.length == 65);

        // We need to unpack the signature, which is given as an array of 65 bytes (like eth.sign)
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(_signature, 32))
            s := mload(add(_signature, 64))
            v := and(mload(add(_signature, 65)), 255)
        }
        if (v < 27) {
            v += 27; // Ethereum versions are 27 or 28 as opposed to 0 or 1 which is submitted by some signing libs
        }
        
        return ecrecover(_operationHash, v, r, s);
    }

    /**
     * @dev Check the number of signatures.
     * @param _operation operation hash
     * @param _signature the multi-signature administrator's signature
     * @return if the number of signatures needs to be collectted is enough, return true.
     */
    function _checkSigNum(bytes32 _operation, bytes memory _signature) private returns (bool) {
        PendingState storage pending = _m_pending[_operation];
        // if we're not yet working on this operation, switch over and reset the confirmation status.

        if (pending.yetNeeded == 0) {
            // reset count of confirmations needed.
            pending.yetNeeded = _m_required;
            // reset which owners have confirmed (none) - set our bitmap to 0.
            pending.ownersDone = 0;
            pending.index = pending.index + 1;
            delete pending.ownersAddr;
        } 

        address sigAddr = _recoverAddressFromSignature(_operation, _signature);

        if (pending.ownersAddr.length != 0) {
            for (uint256 i = 0; i < _m_required; i++) {
                if (pending.ownersAddr[i] == sigAddr) {
                    return false;
                }
            }
        }

        if (pending.yetNeeded <= 1) {
            pending.yetNeeded--;
            pending.ownersDone++;
            return true;
        } else {
            pending.ownersAddr[pending.ownersDone] = sigAddr; 
            pending.yetNeeded--;
            pending.ownersDone++;
            return false;
        }

        return false;
    }

    /**
     * @dev Check the number of signatures which still needs.
     * @param _contractAddr contract's address
     * @param _from transfer's address
     * @param _to recipient's address
     * @param _value the amount of tokens to be spent
     * @return the number of signatures which still needs
     */
    function checkSigNeeded(address _contractAddr, address _from, address _to, uint256 _value) public view returns(uint256) {
        bytes32 operation = keccak256(
            abi.encodePacked(_contractAddr, _from, _to, _value)
        );
        
        return _m_pending[operation].yetNeeded;
    }

    /**
     * @dev Change the administrator.
     * @param _user new administrator
     */
    function changeSigAdmin(address _user) onlyAdmin public{
        _admin = _user;

        changeOwner(msg.sender, _user);

        emit ChangeAdmin(_user);
    }
}

