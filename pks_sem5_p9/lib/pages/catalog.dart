import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '/components/card_preview.dart';
import '/main.dart';
import '/models/shop_item.dart';
import '/pages/add_item.dart';
import '/pages/item_view.dart';


class Catalog extends StatefulWidget {
  const Catalog({super.key});
  @override
  createState() => CatalogState();
}

class CatalogState extends State<Catalog>
{
  List<ShopItem> shopItems = appData.shopItems;
  @override void initState() {
    super.initState();
  }

  void addItem(ShopItem item) {
    setState(() {
      shopItems.add(item);
    });
  }

  void removeItem(int index)
  {
    setState(() {
      int id = shopItems[index].ID;
      shopItems.removeAt(index);
    });
  }
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15),
        child: GridView.builder(
           gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                childAspectRatio: 21/20),
          padding: const EdgeInsets.symmetric(vertical: 0),
          itemCount: shopItems.length,
          itemBuilder: (BuildContext context, int index) {
            return GestureDetector(
              child: CardPreview(shopItem: shopItems[index],),
              onTap: () {
                debugPrint('tapped ${shopItems[index].Name}');
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ItemView(shopItem: shopItems[index]))
                );
              },
            );
          },
        ),
      ),
      floatingActionButton: FloatingActionButton(
        tooltip: 'Add Note',
        onPressed: () {
          Navigator.push(context, MaterialPageRoute(builder: (context) => AddItem(catalogState: this,)));
        },
        child: Icon(Icons.add),
      ),
    );
  }

}
