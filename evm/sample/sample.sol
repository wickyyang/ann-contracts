pragma solidity ^0.5.0;

contract EVMTest {
    struct TaskInfo {
        string id;
        string name;
        string phoneNo;
    }
    
    mapping (string => TaskInfo) TaskInfoIndex;
    
    function call(string memory _id, string memory _name, string memory _phoneNo, uint256 _num) public view returns (string memory, string memory, string memory, uint256, string memory, address) {
        uint256 res = _num + 1;
        string memory str = "you are right";
        
        return (_id, _name, _phoneNo, res, str, msg.sender);
    }
}
