// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

contract Msig {

    struct OperationState {
        uint yetNeeded;
        uint ownersDone;
        uint index;
    }

    uint public m_threshold;
    uint public m_numOwners;

    uint[256] m_owners;

    mapping (uint => uint) m_ownerIndex;
    mapping(bytes32 => OperationState) m_pending;

    mapping(uint => bytes32) m_pendingIndex;
    uint m_pendingIndex_length=0;

    uint totalAmount;


    event Deposit (address accountFrom, uint amount);
    event Withdraw (address accountTo, uint amount);

    event OwnerAdded(address newOwner);
    event OwnerRemoved(address oldOwner);
    event ThresholdChanged(uint newRequirement);
    event Confirmation(address owner, bytes32 operation);


    modifier confirmed(bytes32 operation) {
        if (confirmAndCheck(operation))
            _;
    }


    constructor(address[] memory owners, uint threshold) public {

        m_numOwners = owners.length+1; //первый учатсник
        m_owners[1] = uint(msg.sender);

        m_ownerIndex[uint(msg.sender)] = 1;
        
        for (uint i = 0; i < owners.length; ++i)
        {
            m_owners[2 + i] = uint(owners[i]);
            m_ownerIndex[uint(owners[i])] = 2 + i;
        }
        m_threshold = threshold;
    }


    function depositEth() public payable {

        require (m_ownerIndex[uint(msg.sender)] != 0);
        emit Deposit (msg.sender, msg.value);
    }


    function withdrawEth(address payable accountTo, uint amount) confirmed(keccak256(abi.encodePacked(msg.data, block.number))) public {

        require(address(this).balance >= amount);
        accountTo.transfer(amount);
        emit Withdraw (accountTo, amount);
    }


    function getBalance() public view returns(uint) {
        return address(this).balance;
    }



    function addOwner(address owner) public {
        
        if (m_ownerIndex[uint(owner)] > 0) return;

        clearPending();

        if (m_numOwners >= 250) reorganizeOwners();
        if (m_numOwners >= 250) return;
        
        m_numOwners++; //добавили участника
        m_owners[m_numOwners] = uint(owner);
        m_ownerIndex[uint(owner)] = m_numOwners;
        
        emit OwnerAdded(owner);
    }


    function removeOwner(address owner) public {

        uint ownerIndex = m_ownerIndex[uint(owner)];

        if (ownerIndex == 0 || m_threshold > m_numOwners - 1) return;

        m_owners[ownerIndex] = 0;
        m_ownerIndex[uint(owner)] = 0;
        clearPending();
        reorganizeOwners();
        
        emit OwnerRemoved(owner);
    }


    function changeThreshold(uint newThreshold) public {

        if (newThreshold > m_numOwners || newThreshold < 0) return; //не может превышать число участников
        
        m_threshold = newThreshold;
        clearPending();

        emit ThresholdChanged(newThreshold);
    }


    function confirmAndCheck(bytes32 operation) public returns (bool) { //получить подтверждение и проверить, достаточно ли их

        uint ownerIndex = m_ownerIndex[uint(msg.sender)];
        if (ownerIndex == 0) return false;

        OperationState memory pending = m_pending[operation]; //статус операции
        
        if (pending.yetNeeded == 0) { //если операция новая 
            pending.yetNeeded = m_threshold;
            pending.ownersDone = 0;
            pending.index = m_pendingIndex_length+1;

            m_pendingIndex[pending.index] = operation; //добавляем её в список
            m_pendingIndex_length++;
        }

        uint ownerIndexBit = 2**ownerIndex;

        if (pending.ownersDone & ownerIndexBit == 0) { //подтвеждения от участника еще не было
            emit Confirmation(msg.sender, operation);

            if (pending.yetNeeded <= 1) { //если достаточно
                delete m_pendingIndex[m_pending[operation].index]; //убираем операцию из списка
                delete m_pending[operation];
                m_pendingIndex_length--;
                return true;
            }
            else //если не достаточно
            {
                pending.yetNeeded--;
                pending.ownersDone |= ownerIndexBit;
                return false;
            }
        }
        return false;
    }


    function clearPending() internal { //удаляем ненужные  операции
        uint length = m_pendingIndex_length;
        
        for (uint i = 0; i < length; ++i)
            if (m_pendingIndex[i] != 0)
                delete m_pending[m_pendingIndex[i]];
    }


    function reorganizeOwners() private returns (bool) {
        uint free = 1;
        while (free < m_numOwners)
        {
            while (free < m_numOwners && m_owners[free] != 0) free++;
            while (m_numOwners > 1 && m_owners[m_numOwners] == 0) m_numOwners--;
            if (free < m_numOwners && m_owners[m_numOwners] != 0 && m_owners[free] == 0)
            {
                m_owners[free] = m_owners[m_numOwners];
                m_ownerIndex[m_owners[free]] = free;
                m_owners[m_numOwners] = 0;
            }
        }
    }
}