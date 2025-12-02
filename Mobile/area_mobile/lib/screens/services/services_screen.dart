import 'package:flutter/material.dart';
import '../../services/service_connector.dart';

class ServicesScreen extends StatefulWidget {
  const ServicesScreen({super.key});

  @override
  State<ServicesScreen> createState() => _ServicesScreenState();
}

class _ServicesScreenState extends State<ServicesScreen> {
  bool _loading = false;
  String _status = 'Not connected';
  final ServiceConnector connector = ServiceConnector();

  Future<void> _connectGoogle() async {
    setState(() {
      _loading = true;
      _status = 'Connecting to Google...';
    });

    try {
      await connector.connectToService('google', 'JWT_TOKEN');
      setState(() {
        _status = 'Successfully connected to Google service';
      });
    } catch (e) {
      setState(() {
        _status = 'Error: $e';
      });
    } finally {
      setState(() {
        _loading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Services')),
      body: Center(
        child: _loading
            ? const CircularProgressIndicator()
            : Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    _status,
                    style: const TextStyle(fontSize: 18),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton(
                    onPressed: _connectGoogle,
                    child: const Text('Connect Google Service'),
                  ),
                ],
              ),
      ),
    );
  }
}