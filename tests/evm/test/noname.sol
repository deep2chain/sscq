pragma solidity ^0.4.18;

contract TraceAbility {
    
	address root;
	uint storedData;
	address [] producerGroups;//被授权生产者组
	//address [] retailerGroups； //被授权零售商组
	address [] logisticsGroups;//被授权物流组
    address public lastsend;
	
	struct Commodity{
        string commodityID;//产品号
		string transactionID;//交易号
        string commodityName;//产品名
        //uint produceTime;
        //string producerName;
		uint32 weight;//重量
		
        uint[] timestamps;//0地址表示生产者
		string[] location;//0地址表示生产者
        string[] retailerNames;//0地址表示生产者
        //uint sellTime;
        //string customerName;
        bool isEnding;
		bool exist;
        address owner;
    }
	
	mapping(string => Commodity) commodityMap;
	mapping(address => string) producerName;
	mapping(address => string) logisticsName;
	
	
	//访问者判定修饰
	modifier onlyOwner {
           require(msg.sender == root,"sender not authered");
         _;
    }
	//构造函数
	constructor() 
	public 
	{    
        root = msg.sender;
    }
	//判断是否是被授权的生产商
	function isAutheredProducer() 
	private
	view
	returns(bool)
	{
	   for(uint x = 0; x < producerGroups.length; x++)
	   {
			if(msg.sender == producerGroups[x])
			{
			  return true;
			}
	   }
		return false;
	
	}
	//查询是否是被授权的生产商
	function hasAutheredProducer(address paddress) 
	public
	view
	returns(bool)
	{
	   for(uint x = 0; x < producerGroups.length; x++)
	   {
			if(paddress == producerGroups[x])
			{
			  return true;
			}
	   }
		return false;
	
	}
	//判断是否是被授权的物流
	function isAutheredlogistics() 
	private
	view
	returns(bool)
	{
	   for(uint x = 0; x < logisticsGroups.length; x++)
	   {
			if(msg.sender == logisticsGroups[x])
			{
			  return true;
			}
	   }
		return false;
	
	}
	//查询是否是被授权的物流
	function hasAutheredlogistics(address paddress) 
	public
	view
	returns(bool)
	{
	   for(uint x = 0; x < logisticsGroups.length; x++)
	   {
			if(paddress == logisticsGroups[x])
			{
			  return true;
			}
	   }
		return false;
	
	}
	//增加授权生产者
	function addProducer(address newauther,string name)
	public
	onlyOwner
	returns(bool)
	{
	
		if(hasAutheredProducer(newauther))
		{
		
		}
		else
		{
		    producerGroups.push(newauther);
			producerName[newauther] = name;
		}
	  
			
		
	    
	}
	
	//增加授权物流
	function addlogistics(address newauther,string name)
	public
	onlyOwner
	returns(bool)
	{   
		if(hasAutheredlogistics(newauther))
		{
		
		}
		else
		{
		    logisticsGroups.push(newauther);
			logisticsName[newauther] = name;
		}
	
			
			
		
	    
	}
	
	//产品打包
	function PackagenewCommodity(string ID,string Name,uint Time,uint32 weight,string Location,string tranID)
	external
	{   
	    if(isAutheredProducer())
		{
			if(commodityMap[ID].exist == true)
			{
				return ;
			}
		
		commodityMap[ID].commodityID = ID;
		commodityMap[ID].commodityName = Name;
		commodityMap[ID].weight = weight;
		commodityMap[ID].timestamps.push(Time);
		commodityMap[ID].location.push(Location);
		commodityMap[ID].retailerNames.push(producerName[msg.sender]);
		commodityMap[ID].transactionID = tranID;
		commodityMap[ID].exist = true;
		
		}
		
	}
	//运输状态更新
	function TransferStateUpdate(string ID,uint Time,string Location,bool state)
	external
	{
		if(isAutheredlogistics())
		{
			if(commodityMap[ID].exist == false)
			{
				return ;
			}
		
		commodityMap[ID].timestamps.push(Time);
		commodityMap[ID].location.push(Location);
		commodityMap[ID].retailerNames.push(logisticsName[msg.sender]);
			if(state == true)
			{
				commodityMap[ID].isEnding = true;
			}
		
		}
	}
	
	function getCommodityInfo(string ID)
	external
	view
	//      产品名               重量            生产地                 生产商名称            生产时间     
	returns(string commodityname,uint32 weight,string producelocation,string producename,uint producetime)
	{
	   
	   return (commodityMap[ID].commodityName,commodityMap[ID].weight,(commodityMap[ID].location)[0],(commodityMap[ID].retailerNames)[0],(commodityMap[ID].timestamps)[0]);
	}
	
	function getCommodityInfo2(string ID)
	external
	view
	//      订单状态     订单ID     物流信息长度
	returns(bool end, string tranID,uint length)
	{
	   
	   return (commodityMap[ID].isEnding,commodityMap[ID].transactionID,commodityMap[ID].location.length);
	}
	
	function getCommodityTransferinfo(string ID,uint locate)
	external
	view
	returns(uint time,string location,string name)
	{
	  //if(locate <= commodityMap[ID].location.length - 1)
	  //if(locate <= commodityMap[ID].location.length - 1)
	  {
		return ((commodityMap[ID].timestamps)[locate],(commodityMap[ID].location)[locate],(commodityMap[ID].retailerNames)[locate]);
	  }
	  //else
	 // {
		//return (0,"err1","err2");
	  //}
	 
	}
	

}