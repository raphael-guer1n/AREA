import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import '../../services/service_connector.dart';

class ServicesScreen extends StatefulWidget {
  const ServicesScreen({super.key});

  @override
  State<ServicesScreen> createState() => _ServicesScreenState();
}

class _ServicesScreenState extends State<ServicesScreen> {
  bool _loading = false;
  String _status = 'Not connected';
  late final ServiceConnector _connector;

  @override
  void initState() {
    super.initState();
    _connector = ServiceConnector(baseUrl: dotenv.env['BASE_URL'] ?? '');
  }

  Future<void> _connectGoogle() async {
    setState(() {
      _loading = true;
      _status = 'Connecting to Google...';
    });

    try {
      await _connector.connectToService('google', 'FAKE_JWT_TOKEN');
      setState(() {
        _status = 'Connected to Google successfully!';
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
                  Text(_status, textAlign: TextAlign.center, style: const TextStyle(fontSize: 18)),
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