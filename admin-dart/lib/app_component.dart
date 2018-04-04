import 'package:angular/angular.dart';
import 'package:angular_components/angular_components.dart';
import 'src/registry/registry_component.dart';
import 'src/serviceTable/service_table_component.dart';
import 'src/registry/registry_service.dart';

@Component(
  selector: 'app',
  templateUrl: "app_component.html",
  styleUrls: const ['app_component.css'],
  directives: const [
    materialDirectives,
    RegistryComponent,
    ServiceTableComponent
  ],
  providers: const [materialProviders, RegistryService]
)
class AppComponent {}