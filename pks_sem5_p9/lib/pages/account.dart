import 'package:flutter/material.dart';
import '/main.dart';
import '/pages/account_update.dart';
class AccountPage extends StatefulWidget{
  const AccountPage({super.key});

  @override
  State<AccountPage> createState() => AccountPageState();
}

class AccountPageState extends State<AccountPage> {
  @override
  void initState() {
    super.initState();
    appData.accountPageState = this;
  }
  void forceUpdateState()
  {
    setState(() {
      
    });
  }
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        width: double.infinity,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            SizedBox(height: 100,),
            Icon(Icons.account_circle, size: 200.0),
            SizedBox(height: 10,),
            Text(appData.account!.Name, style: TextStyle(
              fontSize: 22
            ),),
            SizedBox(height: 10,),
            Text(appData.account!.Email),
            SizedBox(height: 10,),
            Text(appData.account!.PhoneNumber),
            SizedBox(height: 20,),
            TextButton(onPressed: (){
              
              Navigator.push(context, MaterialPageRoute(builder: (context) => AccountUpdatePage()));
            }, child: Text("Обновить данные", style: TextStyle(fontSize: 18),),),
          ],
        ),
      ),
      persistentFooterButtons: [
        SizedBox.expand(child: TextButton(onPressed: (){}, child: Text("Выйти", style: TextStyle(fontSize: 18),),))
      ],
    );
  }
}