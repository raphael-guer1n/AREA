import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/area_provider.dart';
import '../../models/area_model.dart';
import '../../theme/colors.dart';

class AreaScreen extends StatefulWidget {
  const AreaScreen({super.key});

  @override
  State<AreaScreen> createState() => _AreaScreenState();
}

class _AreaScreenState extends State<AreaScreen> {
  final _formKey = GlobalKey<FormState>();
  final _summaryController = TextEditingController();
  final _descriptionController = TextEditingController();

  DateTime? _startTime;
  DateTime? _endTime;
  int _delay = 0;

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final provider = context.read<AreaProvider>();

    if (_startTime == null || _endTime == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select start and end times')),
      );
      return;
    }

    final req = CreateEventRequest(
      delay: _delay,
      event: EventModel(
        startTime: _startTime!.toUtc().toIso8601String(),
        endTime: _endTime!.toUtc().toIso8601String(),
        summary: _summaryController.text.trim(),
        description: _descriptionController.text.trim(),
      ),
    );

    await provider.createEvent(req);

    if (provider.error != null) {
      ScaffoldMessenger.of(context)
          .showSnackBar(SnackBar(content: Text(provider.error!)));
    } else {
      ScaffoldMessenger.of(context)
          .showSnackBar(SnackBar(content: Text(provider.message ?? 'Success')));
    }
  }

  Future<void> _pickDateTime(BuildContext context, bool isStart) async {
    final now = DateTime.now();
    final pickedDate = await showDatePicker(
      context: context,
      initialDate: now,
      firstDate: now.subtract(const Duration(days: 1)),
      lastDate: DateTime(2100),
    );
    if (pickedDate == null) return;
    final pickedTime =
        await showTimePicker(context: context, initialTime: TimeOfDay.now());
    if (pickedTime == null) return;

    final dt = DateTime(
      pickedDate.year,
      pickedDate.month,
      pickedDate.day,
      pickedTime.hour,
      pickedTime.minute,
    );

    setState(() {
      if (isStart) {
        _startTime = dt;
      } else {
        _endTime = dt;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<AreaProvider>();

    return Scaffold(
      appBar: AppBar(title: const Text('Create AREA')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              const Text(
                'Action: Timer',
                style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
              ),
              const SizedBox(height: 12),
              TextFormField(
                initialValue: '0',
                decoration: const InputDecoration(labelText: 'Delay (seconds)'),
                keyboardType: TextInputType.number,
                onChanged: (v) {
                  _delay = int.tryParse(v) ?? 0;
                },
              ),
              const SizedBox(height: 24),
              const Text(
                'Reaction: Google Calendar Event',
                style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _summaryController,
                decoration: const InputDecoration(labelText: 'Title'),
                validator: (v) =>
                    v == null || v.isEmpty ? 'Required field' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _descriptionController,
                decoration: const InputDecoration(labelText: 'Description'),
              ),
              const SizedBox(height: 12),
              ListTile(
                title: Text(
                  _startTime == null
                      ? 'Select start time'
                      : 'Start: ${_startTime.toString()}',
                ),
                trailing: const Icon(Icons.calendar_today),
                onTap: () => _pickDateTime(context, true),
              ),
              ListTile(
                title: Text(
                  _endTime == null
                      ? 'Select end time'
                      : 'End: ${_endTime.toString()}',
                ),
                trailing: const Icon(Icons.calendar_today),
                onTap: () => _pickDateTime(context, false),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: provider.isLoading ? null : _submit,
                child: provider.isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('Create AREA'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}