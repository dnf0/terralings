package manifest

import (
	"sync"

	"github.com/dnf0/terralings/internal/models"
)

var (
	manifestInstance *models.Manifest
	once             sync.Once
)

func buildManifest() *models.Manifest {
	return &models.Manifest{
		Chapters: []models.Chapter{
			{
				Number:      1,
				Name:        "01_primitives",
				Title:       "HCL Foundations & Core Primitives",
				Description: "Blocks, attributes, provider requirements and first resources",
				Exercises: []models.Exercise{
					{
						Name:        "primitives01",
						Title:       "Terraform Configuration Block",
						Path:        "exercises/01_primitives/primitives01.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Use required_version = \">= 1.6.0\"",
							"Specify required_providers with local source \"hashicorp/local\"",
						},
						Mode: models.ModeValidate,
					},
					{
						Name:        "primitives02",
						Title:       "First Resource Declaration",
						Path:        "exercises/01_primitives/primitives02.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Declare resource local_file with filename and content",
							"Set filename = \"${path.module}/welcome.txt\" and content = \"Welcome to Terralings!\"",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "primitives03",
						Title:       "Resource Dependencies",
						Path:        "exercises/01_primitives/primitives03.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Reference local_file.first.content in the second resource",
							"Implicit dependencies are created by referencing attributes of other resources",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "primitives04",
						Title:       "String Interpolation & Heredoc",
						Path:        "exercises/01_primitives/primitives04.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Use <<-EOT for indented heredoc multi-line strings",
							"Interpolate variables or expressions inside ${...}",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "primitives05",
						Title:       "Syntax & Formatting",
						Path:        "exercises/01_primitives/primitives05.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Align equals signs and use double quotes for strings",
							"Run tofu fmt or terraform fmt to format your code canonically",
						},
						Mode: models.ModeValidate,
					},
					{
						Name:        "primitives06",
						Title:       "Lifecycle Mechanics",
						Path:        "exercises/01_primitives/primitives06.tf",
						ChapterName: "01_primitives",
						Hints: []string{
							"Ensure terraform_data resource has input and triggers_replace",
							"Use triggers_replace to force recreation on value change",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      2,
				Name:        "02_variables",
				Title:       "Input Variables, Types & Validations",
				Description: "Primitive types, collections, structural objects, and custom validation blocks",
				Exercises: []models.Exercise{
					{
						Name:        "variables01",
						Title:       "Primitive Variable Declarations",
						Path:        "exercises/02_variables/variables01.tf",
						ChapterName: "02_variables",
						Hints: []string{
							"Define type = string, number, or bool for input variables",
							"Provide descriptive description attributes",
						},
						Mode: models.ModeValidate,
					},
					{
						Name:        "variables02",
						Title:       "Collection Types",
						Path:        "exercises/02_variables/variables02.tf",
						ChapterName: "02_variables",
						Hints: []string{
							"Use list(string), map(string), or set(string) type constraints",
							"Index lists with [0] and map keys with [\"key\"] or .key",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "variables03",
						Title:       "Structural Types",
						Path:        "exercises/02_variables/variables03.tf",
						ChapterName: "02_variables",
						Hints: []string{
							"Define object({ name = string, port = number, enabled = optional(bool, true) })",
							"Use tuple([...]) for fixed-length ordered collections",
						},
						Mode: models.ModeValidate,
					},
					{
						Name:        "variables04",
						Title:       "Defaults and Nullable",
						Path:        "exercises/02_variables/variables04.tf",
						ChapterName: "02_variables",
						Hints: []string{
							"Set default values on variable declarations",
							"Use nullable = false to prevent null assignments",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "variables05",
						Title:       "Custom Variable Validations",
						Path:        "exercises/02_variables/variables05.tf",
						ChapterName: "02_variables",
						Hints: []string{
							"Add validation { condition = ... error_message = ... } blocks",
							"The condition expression must refer only to the variable itself via var.<name>",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      3,
				Name:        "03_outputs_locals",
				Title:       "Outputs, Locals & Expressions",
				Description: "Output declarations, sensitive redaction, locals DRY calculations, and expressions",
				Exercises: []models.Exercise{
					{
						Name:        "outputs01",
						Title:       "Defining Outputs & Sensitive Redaction",
						Path:        "exercises/03_outputs_locals/outputs01.tf",
						ChapterName: "03_outputs_locals",
						Hints: []string{
							"Declare output blocks with value, description, and sensitive = true",
							"Sensitive outputs are masked in CLI plans and logs",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "locals01",
						Title:       "Locals for Intermediate Calculations",
						Path:        "exercises/03_outputs_locals/locals01.tf",
						ChapterName: "03_outputs_locals",
						Hints: []string{
							"Use locals { ... } blocks to compute reusable expressions",
							"Reference local values with local.<name>",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "expr01",
						Title:       "Ternary Conditional Expressions",
						Path:        "exercises/03_outputs_locals/expr01.tf",
						ChapterName: "03_outputs_locals",
						Hints: []string{
							"Use condition ? true_val : false_val syntax",
							"Both branches must produce compatible types",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "expr02",
						Title:       "Splat Expressions",
						Path:        "exercises/03_outputs_locals/expr02.tf",
						ChapterName: "03_outputs_locals",
						Hints: []string{
							"Use [*] to extract attributes across list of objects",
							"For example: local.servers[*].ip_address",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      4,
				Name:        "04_functions",
				Title:       "Built-in Functions & Collections",
				Description: "String manipulation, collection operations, encodings, filesystem, and safe evaluation",
				Exercises: []models.Exercise{
					{
						Name:        "func01",
						Title:       "String Manipulation Functions",
						Path:        "exercises/04_functions/func01.tf",
						ChapterName: "04_functions",
						Hints: []string{
							"Use format(), join(), split(), replace(), and lower()",
							"format(\"%s-%d\", var.prefix, var.count)",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "func02",
						Title:       "Collection Operations",
						Path:        "exercises/04_functions/func02.tf",
						ChapterName: "04_functions",
						Hints: []string{
							"Use merge(), lookup(), distinct(), slice(), flatten(), zipmap()",
							"flatten() unwraps nested lists into a single flat list",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "func03",
						Title:       "Data Encodings",
						Path:        "exercises/04_functions/func03.tf",
						ChapterName: "04_functions",
						Hints: []string{
							"Use jsonencode(), yamlencode(), base64encode(), jsondecode()",
							"jsonencode transforms HCL structures into valid JSON strings",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "func04",
						Title:       "Filesystem Functions",
						Path:        "exercises/04_functions/func04.tf",
						ChapterName: "04_functions",
						Hints: []string{
							"Use file(), templatefile(), fileset(), fileexists()",
							"templatefile(path, vars) renders template files with variable substitution",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "func05",
						Title:       "Safe Evaluation Expressions",
						Path:        "exercises/04_functions/func05.tf",
						ChapterName: "04_functions",
						Hints: []string{
							"Use try() and can() to handle potentially failing expressions safely",
							"try(eval, fallback) returns fallback if eval errors",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      5,
				Name:        "05_meta_arguments",
				Title:       "Meta-Arguments & Resource Scaling",
				Description: "Scaling resources with count, for_each, explicit dependencies, and lifecycles",
				Exercises: []models.Exercise{
					{
						Name:        "meta01",
						Title:       "Scaling Resources with Count",
						Path:        "exercises/05_meta_arguments/meta01.tf",
						ChapterName: "05_meta_arguments",
						Hints: []string{
							"Use count = <number> and index resources with count.index",
							"Reference specific instances with <resource>.<name>[index]",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "meta02",
						Title:       "Idempotent Mapping with For Each",
						Path:        "exercises/05_meta_arguments/meta02.tf",
						ChapterName: "05_meta_arguments",
						Hints: []string{
							"Use for_each = toset(...) or a map",
							"Reference current key/value with each.key and each.value",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "meta03",
						Title:       "Explicit Dependency Ordering",
						Path:        "exercises/05_meta_arguments/meta03.tf",
						ChapterName: "05_meta_arguments",
						Hints: []string{
							"Use depends_on = [resource_type.name] for non-attribute dependencies",
							"Only use depends_on when hidden ordering requirements exist",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "meta04",
						Title:       "Resource Lifecycle Blocks",
						Path:        "exercises/05_meta_arguments/meta04.tf",
						ChapterName: "05_meta_arguments",
						Hints: []string{
							"Use lifecycle { create_before_destroy = true } to avoid downtime",
							"Use prevent_destroy = true to guard against accidental destruction",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "meta05",
						Title:       "Dynamic Drift Handling",
						Path:        "exercises/05_meta_arguments/meta05.tf",
						ChapterName: "05_meta_arguments",
						Hints: []string{
							"Use lifecycle { ignore_changes = [tags, ...] }",
							"ignore_changes prevents Terraform from reverting out-of-band updates",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      6,
				Name:        "06_dynamic_blocks",
				Title:       "Dynamic Blocks & Advanced HCL",
				Description: "dynamic blocks iteration, custom iterators, nested dynamic blocks, and conditional emission",
				Exercises: []models.Exercise{
					{
						Name:        "dynamic01",
						Title:       "Basic Dynamic Block Iteration",
						Path:        "exercises/06_dynamic_blocks/dynamic01.tf",
						ChapterName: "06_dynamic_blocks",
						Hints: []string{
							"Use dynamic \"setting\" { for_each = ... content { ... } }",
							"Access current element with setting.value",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "dynamic02",
						Title:       "Dynamic Blocks with Custom Iterator",
						Path:        "exercises/06_dynamic_blocks/dynamic02.tf",
						ChapterName: "06_dynamic_blocks",
						Hints: []string{
							"Use iterator = <name> inside dynamic block",
							"Access current element with <name>.key and <name>.value",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "dynamic03",
						Title:       "Nested Dynamic Blocks",
						Path:        "exercises/06_dynamic_blocks/dynamic03.tf",
						ChapterName: "06_dynamic_blocks",
						Hints: []string{
							"Nest dynamic blocks inside the content block of parent dynamic blocks",
							"Ensure inner iterator references outer block correctly",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "dynamic04",
						Title:       "Conditional Dynamic Block Emission",
						Path:        "exercises/06_dynamic_blocks/dynamic04.tf",
						ChapterName: "06_dynamic_blocks",
						Hints: []string{
							"Pass empty list [] to for_each to emit zero blocks conditionally",
							"Use condition ? [value] : []",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      7,
				Name:        "07_data_sources",
				Title:       "Data Sources & State Querying",
				Description: "Querying local filesystem, archive generation, external JSON sources, and pre/postconditions",
				Exercises: []models.Exercise{
					{
						Name:        "data01",
						Title:       "Local Filesystem Data Sources",
						Path:        "exercises/07_data_sources/data01.tf",
						ChapterName: "07_data_sources",
						Hints: []string{
							"Use data \"local_file\" to read existing files into configuration",
							"Reference data source attributes with data.<type>.<name>.<attr>",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "data02",
						Title:       "Archive File Data Sources",
						Path:        "exercises/07_data_sources/data02.tf",
						ChapterName: "07_data_sources",
						Hints: []string{
							"Use data \"archive_file\" to create zip/tar archives on the fly",
							"Set type = \"zip\", source_dir or source_file, and output_path",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "data03",
						Title:       "External Data Source Queries",
						Path:        "exercises/07_data_sources/data03.tf",
						ChapterName: "07_data_sources",
						Hints: []string{
							"Use data \"external\" to run local scripts returning JSON",
							"Output must be flat string-to-string JSON object",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "data04",
						Title:       "Custom Preconditions and Postconditions",
						Path:        "exercises/07_data_sources/data04.tf",
						ChapterName: "07_data_sources",
						Hints: []string{
							"Add lifecycle { precondition { ... } } or postcondition { ... }",
							"Validate assumptions before or after resource execution",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      8,
				Name:        "08_modules",
				Title:       "Modular Infrastructure Architecture",
				Description: "Child modules, input/output encapsulation, multi-instance deployment, provider passing, and boundaries",
				Exercises: []models.Exercise{
					{
						Name:        "module01",
						Title:       "Building a Clean Child Module",
						Path:        "exercises/08_modules/module01",
						ChapterName: "08_modules",
						Hints: []string{
							"Define input variables in variables.tf and outputs in outputs.tf",
							"Encapsulate implementation details inside the module directory",
						},
						Mode: models.ModeValidate,
					},
					{
						Name:        "module02",
						Title:       "Calling Local Child Modules",
						Path:        "exercises/08_modules/module02",
						ChapterName: "08_modules",
						Hints: []string{
							"Use module \"<name>\" { source = \"./...\" }",
							"Pass required inputs as arguments to the module block",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "module03",
						Title:       "Multi-Instance Module Deployment",
						Path:        "exercises/08_modules/module03",
						ChapterName: "08_modules",
						Hints: []string{
							"Use for_each inside module block to deploy multiple module instances",
							"Access instance outputs via module.<name>[each.key].<output>",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "module04",
						Title:       "Passing Provider Configurations & Aliases",
						Path:        "exercises/08_modules/module04",
						ChapterName: "08_modules",
						Hints: []string{
							"Pass providers = { <child_provider> = <parent_provider> }",
							"Declare configuration_aliases in child module terraform { required_providers { ... } }",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "module05",
						Title:       "Submodule Boundaries & Clean Architecture",
						Path:        "exercises/08_modules/module05",
						ChapterName: "08_modules",
						Hints: []string{
							"Avoid deep submodule nesting antipatterns",
							"Keep module interfaces flat, explicit, and independently testable",
						},
						Mode: models.ModeValidate,
					},
				},
			},
			{
				Number:      9,
				Name:        "09_state_refactoring",
				Title:       "State Management & Refactoring Surgery",
				Description: "Declarative moved blocks, migrating resources to modules, import blocks, and replace triggers",
				Exercises: []models.Exercise{
					{
						Name:        "state01",
						Title:       "Declarative Refactoring with Moved Blocks",
						Path:        "exercises/09_state_refactoring/state01.tf",
						ChapterName: "09_state_refactoring",
						Hints: []string{
							"Use moved { from = <old_addr> to = <new_addr> }",
							"Moved blocks prevent resource destruction during address renames",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "state02",
						Title:       "Migrating Count to For-Each with Moved Blocks",
						Path:        "exercises/09_state_refactoring/state02.tf",
						ChapterName: "09_state_refactoring",
						Hints: []string{
							"Use moved { from = terraform_data.db_cluster[0] to = terraform_data.db_cluster[\"primary\"] }",
							"Convert count = 2 to for_each = toset([\"primary\", \"replica\"])",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "state03",
						Title:       "Declarative Import Blocks",
						Path:        "exercises/09_state_refactoring/state03.tf",
						ChapterName: "09_state_refactoring",
						Hints: []string{
							"Use import { to = <resource>.<name> id = \"<resource_id>\" }",
							"Import blocks onboard existing cloud infrastructure into managed state",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "state04",
						Title:       "Controlled Resource Replacement",
						Path:        "exercises/09_state_refactoring/state04.tf",
						ChapterName: "09_state_refactoring",
						Hints: []string{
							"Use lifecycle { replace_triggered_by = [<address>] }",
							"Force recreation when referenced resources or attributes change",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      10,
				Name:        "10_testing",
				Title:       "Native Unit & Integration Testing (.tftest.hcl)",
				Description: "Test assertions, apply validation, mock providers, and expect_failures",
				Exercises: []models.Exercise{
					{
						Name:        "test01",
						Title:       "Basic Test Assertions with Run Blocks",
						Path:        "exercises/10_testing/test01",
						ChapterName: "10_testing",
						Hints: []string{
							"Define run \"<name>\" { command = plan assert { condition = ... error_message = ... } }",
							"Run assertions against planned attribute values",
						},
						Mode: models.ModeTest,
					},
					{
						Name:        "test02",
						Title:       "Validating Applied Resources in Tests",
						Path:        "exercises/10_testing/test02",
						ChapterName: "10_testing",
						Hints: []string{
							"Use command = apply in run block to test actual provisioned output",
							"Validate state outputs against expected contract values",
						},
						Mode: models.ModeTest,
					},
					{
						Name:        "test03",
						Title:       "Mocking Providers and Resources",
						Path:        "exercises/10_testing/test03",
						ChapterName: "10_testing",
						Hints: []string{
							"Use mock_provider \"local\" { ... } or override_resource in .tftest.hcl",
							"Mock providers without real credentials or network calls",
						},
						Mode: models.ModeTest,
					},
					{
						Name:        "test04",
						Title:       "Testing Failure Cases with Expect Failures",
						Path:        "exercises/10_testing/test04",
						ChapterName: "10_testing",
						Hints: []string{
							"Use expect_failures = [var.<name>] inside run block",
							"Verify that invalid input variables trigger validation errors as expected",
						},
						Mode: models.ModeTest,
					},
				},
			},
			{
				Number:      11,
				Name:        "11_patterns",
				Title:       "Production Patterns & Anti-Patterns",
				Description: "Multi-environment configurations, feature flags, tagging factories, and self-service contracts",
				Exercises: []models.Exercise{
					{
						Name:        "pattern01",
						Title:       "Multi-Environment Configuration Mapping",
						Path:        "exercises/11_patterns/pattern01.tf",
						ChapterName: "11_patterns",
						Hints: []string{
							"Define environment map in locals: local.environments[var.env]",
							"Extract current environment config and pass into terraform_data resource input",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "pattern02",
						Title:       "Feature Flags & Conditional Resource Creation",
						Path:        "exercises/11_patterns/pattern02.tf",
						ChapterName: "11_patterns",
						Hints: []string{
							"Use count = var.enable_disaster_recovery ? 1 : 0 for optional resources",
							"Use one() or try() to safely reference conditional resource outputs",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "pattern03",
						Title:       "Tagging Factory Pattern",
						Path:        "exercises/11_patterns/pattern03.tf",
						ChapterName: "11_patterns",
						Hints: []string{
							"Merge base_tags, env_tags, and custom extra_tags using merge()",
							"Ensure standard tags (ManagedBy, Environment, CostCenter) are consistently applied",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "pattern04",
						Title:       "Self-Service Input Contracts",
						Path:        "exercises/11_patterns/pattern04.tf",
						ChapterName: "11_patterns",
						Hints: []string{
							"Validate complex structured variable contracts with validation blocks",
							"Use for comprehensions with conditions to filter and project service endpoints",
						},
						Mode: models.ModePlan,
					},
				},
			},
			{
				Number:      12,
				Name:        "12_opentofu",
				Title:       "OpenTofu Innovations & Enterprise Features",
				Description: "State encryption at rest, early variable evaluation, and public registry ecosystem",
				Exercises: []models.Exercise{
					{
						Name:        "tofu01",
						Title:       "State Encryption at Rest",
						Path:        "exercises/12_opentofu/tofu01.tf",
						ChapterName: "12_opentofu",
						Hints: []string{
							"Configure terraform { encryption { key_provider \"pbkdf2\" { ... } method \"aes_gcm\" { ... } state { ... } } }",
							"OpenTofu natively encrypts state and plan files at rest",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "tofu02",
						Title:       "Early Variable Evaluation",
						Path:        "exercises/12_opentofu/tofu02.tf",
						ChapterName: "12_opentofu",
						Hints: []string{
							"Reference var.<name> directly within expressions and early-evaluated configurations",
							"OpenTofu evaluates variables before provider and backend initialization",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "tofu03",
						Title:       "OpenTofu Public Registry Integration",
						Path:        "exercises/12_opentofu/tofu03.tf",
						ChapterName: "12_opentofu",
						Hints: []string{
							"Configure required_providers with required version and source",
							"Ensure open provider ecosystem compatibility",
						},
						Mode: models.ModeValidate,
					},
				},
			},
			{
				Number:      13,
				Name:        "13_governance",
				Title:       "Architecture Governance & Enterprise Standards",
				Description: "Root module encapsulation, policy encapsulation (ADR-0005), and ephemeral workload isolation",
				Exercises: []models.Exercise{
					{
						Name:        "gov01",
						Title:       "Root Module Encapsulation",
						Path:        "exercises/13_governance/gov01.tf",
						ChapterName: "13_governance",
						Hints: []string{
							"Zero plain/loose workload compute resources in root environments",
							"Encapsulate workload primitives (compute, log groups, security groups) into dedicated modules",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "gov02",
						Title:       "Policy Encapsulation (ADR-0005)",
						Path:        "exercises/13_governance/gov02.tf",
						ChapterName: "13_governance",
						Hints: []string{
							"The module that owns the resource owns the policies that talk to it",
							"Pass managed policy ARNs directly into IAM role policy_arns instead of writing inline wildcard JSON policies",
						},
						Mode: models.ModePlan,
					},
					{
						Name:        "gov03",
						Title:       "Ephemeral Workload Isolation",
						Path:        "exercises/13_governance/gov03.tf",
						ChapterName: "13_governance",
						Hints: []string{
							"Encapsulate ephemeral tooling and batch workloads into dedicated namespaces/modules",
							"Prevent ephemeral resources from polluting persistent infrastructure root state",
						},
						Mode: models.ModePlan,
					},
				},
			},
		},
	}
}

// GetManifest returns the singleton curriculum manifest instance.
func GetManifest() *models.Manifest {
	once.Do(func() {
		manifestInstance = buildManifest()
	})
	return manifestInstance
}

// GetExerciseByName finds an exercise by its short name or relative file path.
func GetExerciseByName(name string) *models.Exercise {
	for _, ex := range GetManifest().AllExercises() {
		if ex.Name == name || ex.Path == name {
			e := ex
			return &e
		}
	}
	return nil
}

// GetNextExercise returns the next sequential exercise after currentName, or nil if at the end.
func GetNextExercise(currentName string) *models.Exercise {
	all := GetManifest().AllExercises()
	for i, ex := range all {
		if ex.Name == currentName || ex.Path == currentName {
			if i+1 < len(all) {
				next := all[i+1]
				return &next
			}
			return nil
		}
	}
	return nil
}
