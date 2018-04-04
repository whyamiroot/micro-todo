import 'dart:collection';
import 'package:angular/angular.dart';
import 'package:admin_dart/src/service/service_component.dart';
import 'package:admin_dart/src/service/service.dart';
import 'package:admin_dart/src/registry/registry_service.dart';

@Component(
  selector: 'service-table',
  templateUrl: 'service_table_component.html',
  styleUrls: const ['service_table_component.css'],
  directives: const [CORE_DIRECTIVES, ServiceComponent]
)
class ServiceTableComponent {
  final RegistryService registryService;

  ServiceTableComponent(RegistryService this.registryService);

  HashMap<String, List<Service>> services;
}