import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';

export type ExerciseStatus = 'passed' | 'failed' | 'in_progress' | 'not_started';

export interface ExerciseDefinition {
  name: string;
  title: string;
  path: string;
  chapterName: string;
  hints?: string[];
  mode?: string;
}

export interface ChapterDefinition {
  number: number;
  name: string;
  title: string;
  description: string;
  exercises: ExerciseDefinition[];
}

export interface ExerciseStateJson {
  name: string;
  chapter?: string;
  status?: string;
  attempts?: number;
  hints_viewed?: number;
  first_attempt_at?: string;
  completed_at?: string;
  time_spent_seconds?: number;
}

export interface StateDataJson {
  version?: string;
  created_at?: string;
  last_active_at?: string;
  total_time_spent_seconds?: number;
  exercises?: Record<string, ExerciseStateJson>;
}

export const CURRICULUM: ChapterDefinition[] = [
  {
    number: 1,
    name: '01_primitives',
    title: 'HCL Foundations & Core Primitives',
    description: 'Blocks, attributes, provider requirements and first resources',
    exercises: [
      {
        name: 'primitives01',
        title: 'Terraform Configuration Block',
        path: 'exercises/01_primitives/primitives01.tf',
        chapterName: '01_primitives',
        hints: [
          'Use required_version = ">= 1.6.0"',
          'Specify required_providers with local source "hashicorp/local"'
        ],
        mode: 'validate'
      },
      {
        name: 'primitives02',
        title: 'First Resource Declaration',
        path: 'exercises/01_primitives/primitives02.tf',
        chapterName: '01_primitives',
        hints: [
          'Declare resource local_file with filename and content',
          'Set filename = "${path.module}/welcome.txt" and content = "Welcome to Terralings!"'
        ],
        mode: 'plan'
      },
      {
        name: 'primitives03',
        title: 'Resource Dependencies',
        path: 'exercises/01_primitives/primitives03.tf',
        chapterName: '01_primitives',
        hints: [
          'Reference local_file.first.content in the second resource',
          'Implicit dependencies are created by referencing attributes of other resources'
        ],
        mode: 'plan'
      },
      {
        name: 'primitives04',
        title: 'String Interpolation & Heredoc',
        path: 'exercises/01_primitives/primitives04.tf',
        chapterName: '01_primitives',
        hints: [
          'Use <<-EOT for indented heredoc multi-line strings',
          'Interpolate variables or expressions inside ${...}'
        ],
        mode: 'plan'
      },
      {
        name: 'primitives05',
        title: 'Syntax & Formatting',
        path: 'exercises/01_primitives/primitives05.tf',
        chapterName: '01_primitives',
        hints: [
          'Align equals signs and use double quotes for strings',
          'Run tofu fmt or terraform fmt to format your code canonically'
        ],
        mode: 'validate'
      },
      {
        name: 'primitives06',
        title: 'Lifecycle Mechanics',
        path: 'exercises/01_primitives/primitives06.tf',
        chapterName: '01_primitives',
        hints: [
          'Ensure terraform_data resource has input and triggers_replace',
          'Use triggers_replace to force recreation on value change'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 2,
    name: '02_variables',
    title: 'Input Variables, Types & Validations',
    description: 'Primitive types, collections, structural objects, and custom validation blocks',
    exercises: [
      {
        name: 'variables01',
        title: 'Primitive Variable Declarations',
        path: 'exercises/02_variables/variables01.tf',
        chapterName: '02_variables',
        hints: [
          'Define type = string, number, or bool for input variables',
          'Provide descriptive description attributes'
        ],
        mode: 'validate'
      },
      {
        name: 'variables02',
        title: 'Collection Types',
        path: 'exercises/02_variables/variables02.tf',
        chapterName: '02_variables',
        hints: [
          'Use list(string), map(string), or set(string) type constraints',
          'Index lists with [0] and map keys with ["key"] or .key'
        ],
        mode: 'plan'
      },
      {
        name: 'variables03',
        title: 'Structural Types',
        path: 'exercises/02_variables/variables03.tf',
        chapterName: '02_variables',
        hints: [
          'Define object({ name = string, port = number, enabled = optional(bool, true) })',
          'Use tuple([...]) for fixed-length ordered collections'
        ],
        mode: 'validate'
      },
      {
        name: 'variables04',
        title: 'Defaults and Nullable',
        path: 'exercises/02_variables/variables04.tf',
        chapterName: '02_variables',
        hints: [
          'Set default values on variable declarations',
          'Use nullable = false to prevent null assignments'
        ],
        mode: 'plan'
      },
      {
        name: 'variables05',
        title: 'Custom Variable Validations',
        path: 'exercises/02_variables/variables05.tf',
        chapterName: '02_variables',
        hints: [
          'Add validation { condition = ... error_message = ... } blocks',
          'The condition expression must refer only to the variable itself via var.<name>'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 3,
    name: '03_outputs_locals',
    title: 'Outputs, Locals & Expressions',
    description: 'Output declarations, sensitive redaction, locals DRY calculations, and expressions',
    exercises: [
      {
        name: 'outputs01',
        title: 'Defining Outputs & Sensitive Redaction',
        path: 'exercises/03_outputs_locals/outputs01.tf',
        chapterName: '03_outputs_locals',
        hints: [
          'Declare output blocks with value, description, and sensitive = true',
          'Sensitive outputs are masked in CLI plans and logs'
        ],
        mode: 'plan'
      },
      {
        name: 'locals01',
        title: 'Locals for Intermediate Calculations',
        path: 'exercises/03_outputs_locals/locals01.tf',
        chapterName: '03_outputs_locals',
        hints: [
          'Use locals { ... } blocks to compute reusable expressions',
          'Reference local values with local.<name>'
        ],
        mode: 'plan'
      },
      {
        name: 'expr01',
        title: 'Ternary Conditional Expressions',
        path: 'exercises/03_outputs_locals/expr01.tf',
        chapterName: '03_outputs_locals',
        hints: [
          'Use condition ? true_val : false_val syntax',
          'Both branches must produce compatible types'
        ],
        mode: 'plan'
      },
      {
        name: 'expr02',
        title: 'Splat Expressions',
        path: 'exercises/03_outputs_locals/expr02.tf',
        chapterName: '03_outputs_locals',
        hints: [
          'Use [*] to extract attributes across list of objects',
          'For example: local.servers[*].ip_address'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 4,
    name: '04_functions',
    title: 'Built-in Functions & Collections',
    description: 'String manipulation, collection operations, encodings, filesystem, and safe evaluation',
    exercises: [
      {
        name: 'func01',
        title: 'String Manipulation Functions',
        path: 'exercises/04_functions/func01.tf',
        chapterName: '04_functions',
        hints: [
          'Use format(), join(), split(), replace(), and lower()',
          'format("%s-%d", var.prefix, var.count)'
        ],
        mode: 'plan'
      },
      {
        name: 'func02',
        title: 'Collection Operations',
        path: 'exercises/04_functions/func02.tf',
        chapterName: '04_functions',
        hints: [
          'Use merge(), lookup(), distinct(), slice(), flatten(), zipmap()',
          'flatten() unwraps nested lists into a single flat list'
        ],
        mode: 'plan'
      },
      {
        name: 'func03',
        title: 'Data Encodings',
        path: 'exercises/04_functions/func03.tf',
        chapterName: '04_functions',
        hints: [
          'Use jsonencode(), yamlencode(), base64encode(), jsondecode()',
          'jsonencode transforms HCL structures into valid JSON strings'
        ],
        mode: 'plan'
      },
      {
        name: 'func04',
        title: 'Filesystem Functions',
        path: 'exercises/04_functions/func04.tf',
        chapterName: '04_functions',
        hints: [
          'Use file(), templatefile(), fileset(), fileexists()',
          'templatefile(path, vars) renders template files with variable substitution'
        ],
        mode: 'plan'
      },
      {
        name: 'func05',
        title: 'Safe Evaluation Expressions',
        path: 'exercises/04_functions/func05.tf',
        chapterName: '04_functions',
        hints: [
          'Use try() and can() to handle potentially failing expressions safely',
          'try(eval, fallback) returns fallback if eval errors'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 5,
    name: '05_meta_arguments',
    title: 'Meta-Arguments & Resource Scaling',
    description: 'Scaling resources with count, for_each, explicit dependencies, and lifecycles',
    exercises: [
      {
        name: 'meta01',
        title: 'Scaling Resources with Count',
        path: 'exercises/05_meta_arguments/meta01.tf',
        chapterName: '05_meta_arguments',
        hints: [
          'Use count = <number> and index resources with count.index',
          'Reference specific instances with <resource>.<name>[index]'
        ],
        mode: 'plan'
      },
      {
        name: 'meta02',
        title: 'Idempotent Mapping with For Each',
        path: 'exercises/05_meta_arguments/meta02.tf',
        chapterName: '05_meta_arguments',
        hints: [
          'Use for_each = toset(...) or a map',
          'Reference current key/value with each.key and each.value'
        ],
        mode: 'plan'
      },
      {
        name: 'meta03',
        title: 'Explicit Dependency Ordering',
        path: 'exercises/05_meta_arguments/meta03.tf',
        chapterName: '05_meta_arguments',
        hints: [
          'Use depends_on = [resource_type.name] for non-attribute dependencies',
          'Only use depends_on when hidden ordering requirements exist'
        ],
        mode: 'plan'
      },
      {
        name: 'meta04',
        title: 'Resource Lifecycle Blocks',
        path: 'exercises/05_meta_arguments/meta04.tf',
        chapterName: '05_meta_arguments',
        hints: [
          'Use lifecycle { create_before_destroy = true } to avoid downtime',
          'Use prevent_destroy = true to guard against accidental destruction'
        ],
        mode: 'plan'
      },
      {
        name: 'meta05',
        title: 'Dynamic Drift Handling',
        path: 'exercises/05_meta_arguments/meta05.tf',
        chapterName: '05_meta_arguments',
        hints: [
          'Use lifecycle { ignore_changes = [tags, ...] }',
          'ignore_changes prevents Terraform from reverting out-of-band updates'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 6,
    name: '06_dynamic_blocks',
    title: 'Dynamic Blocks & Advanced HCL',
    description: 'dynamic blocks iteration, custom iterators, nested dynamic blocks, and conditional emission',
    exercises: [
      {
        name: 'dynamic01',
        title: 'Basic Dynamic Block Iteration',
        path: 'exercises/06_dynamic_blocks/dynamic01.tf',
        chapterName: '06_dynamic_blocks',
        hints: [
          'Use dynamic "setting" { for_each = ... content { ... } }',
          'Access current element with setting.value'
        ],
        mode: 'plan'
      },
      {
        name: 'dynamic02',
        title: 'Dynamic Blocks with Custom Iterator',
        path: 'exercises/06_dynamic_blocks/dynamic02.tf',
        chapterName: '06_dynamic_blocks',
        hints: [
          'Use iterator = <name> inside dynamic block',
          'Access current element with <name>.key and <name>.value'
        ],
        mode: 'plan'
      },
      {
        name: 'dynamic03',
        title: 'Nested Dynamic Blocks',
        path: 'exercises/06_dynamic_blocks/dynamic03.tf',
        chapterName: '06_dynamic_blocks',
        hints: [
          'Nest dynamic blocks inside the content block of parent dynamic blocks',
          'Ensure inner iterator references outer block correctly'
        ],
        mode: 'plan'
      },
      {
        name: 'dynamic04',
        title: 'Conditional Dynamic Block Emission',
        path: 'exercises/06_dynamic_blocks/dynamic04.tf',
        chapterName: '06_dynamic_blocks',
        hints: [
          'Pass empty list [] to for_each to emit zero blocks conditionally',
          'Use condition ? [value] : []'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 7,
    name: '07_data_sources',
    title: 'Data Sources & State Querying',
    description: 'Querying local filesystem, archive generation, external JSON sources, and pre/postconditions',
    exercises: [
      {
        name: 'data01',
        title: 'Local Filesystem Data Sources',
        path: 'exercises/07_data_sources/data01.tf',
        chapterName: '07_data_sources',
        hints: [
          'Use data "local_file" to read existing files into configuration',
          'Reference data source attributes with data.<type>.<name>.<attr>'
        ],
        mode: 'plan'
      },
      {
        name: 'data02',
        title: 'Archive File Data Sources',
        path: 'exercises/07_data_sources/data02.tf',
        chapterName: '07_data_sources',
        hints: [
          'Use data "archive_file" to create zip/tar archives on the fly',
          'Set type = "zip", source_dir or source_file, and output_path'
        ],
        mode: 'plan'
      },
      {
        name: 'data03',
        title: 'External Data Source Queries',
        path: 'exercises/07_data_sources/data03.tf',
        chapterName: '07_data_sources',
        hints: [
          'Use data "external" to run local scripts returning JSON',
          'Output must be flat string-to-string JSON object'
        ],
        mode: 'plan'
      },
      {
        name: 'data04',
        title: 'Custom Preconditions and Postconditions',
        path: 'exercises/07_data_sources/data04.tf',
        chapterName: '07_data_sources',
        hints: [
          'Add lifecycle { precondition { ... } } or postcondition { ... }',
          'Validate assumptions before or after resource execution'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 8,
    name: '08_modules',
    title: 'Modular Infrastructure Architecture',
    description: 'Child modules, input/output encapsulation, multi-instance deployment, provider passing, and boundaries',
    exercises: [
      {
        name: 'module01',
        title: 'Building a Clean Child Module',
        path: 'exercises/08_modules/module01',
        chapterName: '08_modules',
        hints: [
          'Define input variables in variables.tf and outputs in outputs.tf',
          'Encapsulate implementation details inside the module directory'
        ],
        mode: 'validate'
      },
      {
        name: 'module02',
        title: 'Calling Local Child Modules',
        path: 'exercises/08_modules/module02',
        chapterName: '08_modules',
        hints: [
          'Use module "<name>" { source = "./..." }',
          'Pass required inputs as arguments to the module block'
        ],
        mode: 'plan'
      },
      {
        name: 'module03',
        title: 'Multi-Instance Module Deployment',
        path: 'exercises/08_modules/module03',
        chapterName: '08_modules',
        hints: [
          'Use for_each inside module block to deploy multiple module instances',
          'Access instance outputs via module.<name>[each.key].<output>'
        ],
        mode: 'plan'
      },
      {
        name: 'module04',
        title: 'Passing Provider Configurations & Aliases',
        path: 'exercises/08_modules/module04',
        chapterName: '08_modules',
        hints: [
          'Pass providers = { <child_provider> = <parent_provider> }',
          'Declare configuration_aliases in child module terraform { required_providers { ... } }'
        ],
        mode: 'plan'
      },
      {
        name: 'module05',
        title: 'Submodule Boundaries & Clean Architecture',
        path: 'exercises/08_modules/module05',
        chapterName: '08_modules',
        hints: [
          'Avoid deep submodule nesting antipatterns',
          'Keep module interfaces flat, explicit, and independently testable'
        ],
        mode: 'validate'
      }
    ]
  },
  {
    number: 9,
    name: '09_state_refactoring',
    title: 'State Management & Refactoring Surgery',
    description: 'Declarative moved blocks, migrating resources to modules, import blocks, and replace triggers',
    exercises: [
      {
        name: 'state01',
        title: 'Declarative Refactoring with Moved Blocks',
        path: 'exercises/09_state_refactoring/state01.tf',
        chapterName: '09_state_refactoring',
        hints: [
          'Use moved { from = <old_addr> to = <new_addr> }',
          'Moved blocks prevent resource destruction during address renames'
        ],
        mode: 'plan'
      },
      {
        name: 'state02',
        title: 'Migrating Count to For-Each with Moved Blocks',
        path: 'exercises/09_state_refactoring/state02.tf',
        chapterName: '09_state_refactoring',
        hints: [
          'Use moved { from = terraform_data.db_cluster[0] to = terraform_data.db_cluster["primary"] }',
          'Convert count = 2 to for_each = toset(["primary", "replica"])'
        ],
        mode: 'plan'
      },
      {
        name: 'state03',
        title: 'Declarative Import Blocks',
        path: 'exercises/09_state_refactoring/state03.tf',
        chapterName: '09_state_refactoring',
        hints: [
          'Use import { to = <resource>.<name> id = "<resource_id>" }',
          'Import blocks onboard existing cloud infrastructure into managed state'
        ],
        mode: 'plan'
      },
      {
        name: 'state04',
        title: 'Controlled Resource Replacement',
        path: 'exercises/09_state_refactoring/state04.tf',
        chapterName: '09_state_refactoring',
        hints: [
          'Use lifecycle { replace_triggered_by = [<address>] }',
          'Force recreation when referenced resources or attributes change'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 10,
    name: '10_testing',
    title: 'Native Unit & Integration Testing (.tftest.hcl)',
    description: 'Test assertions, apply validation, mock providers, and expect_failures',
    exercises: [
      {
        name: 'test01',
        title: 'Basic Test Assertions with Run Blocks',
        path: 'exercises/10_testing/test01',
        chapterName: '10_testing',
        hints: [
          'Define run "<name>" { command = plan assert { condition = ... error_message = ... } }',
          'Run assertions against planned attribute values'
        ],
        mode: 'test'
      },
      {
        name: 'test02',
        title: 'Validating Applied Resources in Tests',
        path: 'exercises/10_testing/test02',
        chapterName: '10_testing',
        hints: [
          'Use command = apply in run block to test actual provisioned output',
          'Validate state outputs against expected contract values'
        ],
        mode: 'test'
      },
      {
        name: 'test03',
        title: 'Mocking Providers and Resources',
        path: 'exercises/10_testing/test03',
        chapterName: '10_testing',
        hints: [
          'Use mock_provider "local" { ... } or override_resource in .tftest.hcl',
          'Mock providers without real credentials or network calls'
        ],
        mode: 'test'
      },
      {
        name: 'test04',
        title: 'Testing Failure Cases with Expect Failures',
        path: 'exercises/10_testing/test04',
        chapterName: '10_testing',
        hints: [
          'Use expect_failures = [var.<name>] inside run block',
          'Verify that invalid input variables trigger validation errors as expected'
        ],
        mode: 'test'
      }
    ]
  },
  {
    number: 11,
    name: '11_patterns',
    title: 'Production Patterns & Anti-Patterns',
    description: 'Multi-Environment configurations, feature flags, tagging factories, and self-service contracts',
    exercises: [
      {
        name: 'pattern01',
        title: 'Multi-Environment Configuration Mapping',
        path: 'exercises/11_patterns/pattern01.tf',
        chapterName: '11_patterns',
        hints: [
          'Define environment map in locals: local.environments[var.env]',
          'Extract current environment config and pass into terraform_data resource input'
        ],
        mode: 'plan'
      },
      {
        name: 'pattern02',
        title: 'Feature Flags & Conditional Resource Creation',
        path: 'exercises/11_patterns/pattern02.tf',
        chapterName: '11_patterns',
        hints: [
          'Use count = var.enable_disaster_recovery ? 1 : 0 for optional resources',
          'Use one() or try() to safely reference conditional resource outputs'
        ],
        mode: 'plan'
      },
      {
        name: 'pattern03',
        title: 'Tagging Factory Pattern',
        path: 'exercises/11_patterns/pattern03.tf',
        chapterName: '11_patterns',
        hints: [
          'Merge base_tags, env_tags, and custom extra_tags using merge()',
          'Ensure standard tags (ManagedBy, Environment, CostCenter) are consistently applied'
        ],
        mode: 'plan'
      },
      {
        name: 'pattern04',
        title: 'Self-Service Input Contracts',
        path: 'exercises/11_patterns/pattern04.tf',
        chapterName: '11_patterns',
        hints: [
          'Validate complex structured variable contracts with validation blocks',
          'Use for comprehensions with conditions to filter and project service endpoints'
        ],
        mode: 'plan'
      }
    ]
  },
  {
    number: 12,
    name: '12_opentofu',
    title: 'OpenTofu Innovations & Enterprise Features',
    description: 'State Encryption at Rest, early variable evaluation, and public registry ecosystem',
    exercises: [
      {
        name: 'tofu01',
        title: 'State Encryption at Rest',
        path: 'exercises/12_opentofu/tofu01.tf',
        chapterName: '12_opentofu',
        hints: [
          'Configure terraform { encryption { key_provider "pbkdf2" { ... } method "aes_gcm" { ... } state { ... } } }',
          'OpenTofu natively encrypts state and plan files at rest'
        ],
        mode: 'plan'
      },
      {
        name: 'tofu02',
        title: 'Early Variable Evaluation',
        path: 'exercises/12_opentofu/tofu02.tf',
        chapterName: '12_opentofu',
        hints: [
          'Reference var.<name> directly within expressions and early-evaluated configurations',
          'OpenTofu evaluates variables before provider and backend initialization'
        ],
        mode: 'plan'
      },
      {
        name: 'tofu03',
        title: 'OpenTofu Public Registry Integration',
        path: 'exercises/12_opentofu/tofu03.tf',
        chapterName: '12_opentofu',
        hints: [
          'Configure required_providers with required version and source',
          'Ensure open provider ecosystem compatibility'
        ],
        mode: 'validate'
      }
    ]
  },
  {
    number: 13,
    name: '13_governance',
    title: 'Architecture Governance & Enterprise Standards',
    description: 'Root module encapsulation, policy encapsulation (ADR-0005), and ephemeral workload isolation',
    exercises: [
      {
        name: 'gov01',
        title: 'Root Module Encapsulation',
        path: 'exercises/13_governance/gov01.tf',
        chapterName: '13_governance',
        hints: [
          'Zero plain/loose workload compute resources in root environments',
          'Encapsulate workload primitives (compute, log groups, security groups) into dedicated modules'
        ],
        mode: 'plan'
      },
      {
        name: 'gov02',
        title: 'Policy Encapsulation (ADR-0005)',
        path: 'exercises/13_governance/gov02.tf',
        chapterName: '13_governance',
        hints: [
          'The module that owns the resource owns the policies that talk to it',
          'Pass managed policy ARNs directly into IAM role policy_arns instead of writing inline wildcard JSON policies'
        ],
        mode: 'plan'
      },
      {
        name: 'gov03',
        title: 'Ephemeral Workload Isolation',
        path: 'exercises/13_governance/gov03.tf',
        chapterName: '13_governance',
        hints: [
          'Encapsulate ephemeral tooling and batch workloads into dedicated namespaces/modules',
          'Prevent ephemeral resources from polluting persistent infrastructure root state'
        ],
        mode: 'plan'
      }
    ]
  }
];

export const ALL_EXERCISES: ExerciseDefinition[] = CURRICULUM.flatMap((c) => c.exercises);

/**
 * Finds an exercise definition by its unique name or relative path.
 */
export function findExercise(nameOrPath: string): ExerciseDefinition | undefined {
  const trimmed = nameOrPath.trim();
  return ALL_EXERCISES.find(
    (ex) => ex.name === trimmed || ex.path === trimmed || ex.path.endsWith(trimmed)
  );
}

/**
 * Finds a chapter definition by chapter number or chapter name.
 */
export function findChapter(nameOrNumber: string | number): ChapterDefinition | undefined {
  if (typeof nameOrNumber === 'number') {
    return CURRICULUM.find((c) => c.number === nameOrNumber);
  }
  return CURRICULUM.find((c) => c.name === nameOrNumber || c.title === nameOrNumber);
}

/**
 * Resolves the absolute file or directory path for an exercise.
 */
export function resolveExerciseUri(exercisePath: string, workspaceRoot?: string): vscode.Uri {
  const root = workspaceRoot ?? vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  let fullPath = exercisePath;
  if (root && !path.isAbsolute(fullPath)) {
    fullPath = path.join(root, fullPath);
  }
  if (!fullPath.endsWith('.tf') && !fullPath.endsWith('.hcl')) {
    const mainTf = path.join(fullPath, 'main.tf');
    if (fs.existsSync(mainTf)) {
      fullPath = mainTf;
    }
  }
  return vscode.Uri.file(fullPath);
}

export type TreeItemType = 'chapter' | 'exercise';

/**
 * Custom TreeItem representation for chapters and exercises in the Curriculum view.
 */
export class TerralingsTreeItem extends vscode.TreeItem {
  constructor(
    public readonly itemType: TreeItemType,
    public readonly chapter?: ChapterDefinition,
    public readonly exercise?: ExerciseDefinition,
    public readonly status?: ExerciseStatus,
    public readonly completedCount?: number,
    public readonly totalCount?: number
  ) {
    super(
      itemType === 'chapter'
        ? `${chapter ? `${chapter.number}. ` : ''}${chapter?.title ?? ''}`
        : `${exercise?.name}: ${exercise?.title ?? ''}`,
      itemType === 'chapter'
        ? vscode.TreeItemCollapsibleState.Collapsed
        : vscode.TreeItemCollapsibleState.None
    );

    if (itemType === 'chapter' && chapter) {
      this.contextValue = 'chapter';
      this.id = `chapter-${chapter.name}`;
      const completed = completedCount ?? 0;
      const total = totalCount ?? chapter.exercises.length;
      this.description = `${completed}/${total}`;
      this.tooltip = `${chapter.title}\n${chapter.description}\nProgress: ${completed}/${total} completed`;
      this.iconPath =
        completed === total && total > 0
          ? new vscode.ThemeIcon('folder-active')
          : new vscode.ThemeIcon('folder');
    } else if (itemType === 'exercise' && exercise) {
      this.contextValue = 'exercise';
      this.id = `exercise-${exercise.name}`;

      switch (status) {
        case 'passed':
          this.iconPath = new vscode.ThemeIcon(
            'pass-filled',
            new vscode.ThemeColor('testing.iconPassed')
          );
          this.description = 'Passed';
          break;
        case 'failed':
          this.iconPath = new vscode.ThemeIcon(
            'error',
            new vscode.ThemeColor('testing.iconFailed')
          );
          this.description = 'Failed';
          break;
        case 'in_progress':
          this.iconPath = new vscode.ThemeIcon(
            'play',
            new vscode.ThemeColor('testing.iconQueued')
          );
          this.description = 'In Progress';
          break;
        case 'not_started':
        default:
          this.iconPath = new vscode.ThemeIcon('circle-outline');
          this.description = '';
          break;
      }

      const modeStr = exercise.mode ? ` [mode: ${exercise.mode}]` : '';
      this.tooltip = `${exercise.name}: ${exercise.title}${modeStr}\nPath: ${exercise.path}\nStatus: ${status ?? 'not_started'}`;

      const targetUri = resolveExerciseUri(exercise.path);
      this.command = {
        command: 'vscode.open',
        title: 'Open Exercise',
        arguments: [targetUri]
      };
    }
  }
}

/**
 * TreeDataProvider driving the Curriculum & Exercises view in the Terralings sidebar.
 */
export class TerralingsTreeDataProvider
  implements vscode.TreeDataProvider<TerralingsTreeItem>, vscode.Disposable {
  private readonly _onDidChangeTreeData: vscode.EventEmitter<
    TerralingsTreeItem | undefined | null | void
  > = new vscode.EventEmitter<TerralingsTreeItem | undefined | null | void>();
  public readonly onDidChangeTreeData: vscode.Event<
    TerralingsTreeItem | undefined | null | void
  > = this._onDidChangeTreeData.event;

  private exerciseStates: Map<string, ExerciseStatus> = new Map();

  constructor() {
    this.refresh();
  }

  public getTreeItem(element: TerralingsTreeItem): vscode.TreeItem {
    return element;
  }

  public getChildren(element?: TerralingsTreeItem): vscode.ProviderResult<TerralingsTreeItem[]> {
    if (!element) {
      return CURRICULUM.map((chapter) => {
        const total = chapter.exercises.length;
        const completed = chapter.exercises.filter(
          (ex) => this.exerciseStates.get(ex.name) === 'passed'
        ).length;
        return new TerralingsTreeItem(
          'chapter',
          chapter,
          undefined,
          undefined,
          completed,
          total
        );
      });
    }

    if (element.itemType === 'chapter' && element.chapter) {
      return element.chapter.exercises.map((ex) => {
        const status = this.exerciseStates.get(ex.name) ?? 'not_started';
        return new TerralingsTreeItem('exercise', element.chapter, ex, status);
      });
    }

    return [];
  }

  /**
   * Reads .terralings/state.json from the workspace, updates in-memory status map,
   * and fires the change event.
   */
  public refresh(): void {
    this.exerciseStates.clear();
    const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;

    if (workspaceRoot) {
      const stateFile = path.join(workspaceRoot, '.terralings', 'state.json');
      try {
        if (fs.existsSync(stateFile)) {
          const raw = fs.readFileSync(stateFile, 'utf-8');
          const data: StateDataJson = JSON.parse(raw);
          if (data && data.exercises) {
            for (const [name, exState] of Object.entries(data.exercises)) {
              if (!exState) {
                continue;
              }
              if (exState.status === 'passed') {
                this.exerciseStates.set(name, 'passed');
              } else if (exState.status === 'failed') {
                this.exerciseStates.set(name, 'failed');
              } else if (
                exState.status === 'in_progress' ||
                (typeof exState.attempts === 'number' && exState.attempts > 0)
              ) {
                this.exerciseStates.set(name, 'in_progress');
              } else {
                this.exerciseStates.set(name, 'not_started');
              }
            }
          }
        }
      } catch (err) {
        console.error('Failed to parse .terralings/state.json:', err);
      }
    }

    this._onDidChangeTreeData.fire();
  }

  /**
   * Returns progress statistics across the entire curriculum.
   */
  public getProgress(): { completed: number; total: number; percentage: number } {
    const total = ALL_EXERCISES.length;
    let completed = 0;
    for (const ex of ALL_EXERCISES) {
      if (this.exerciseStates.get(ex.name) === 'passed') {
        completed++;
      }
    }
    const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
    return { completed, total, percentage };
  }

  /**
   * Returns current status of a specific exercise.
   */
  public getExerciseStatus(name: string): ExerciseStatus {
    return this.exerciseStates.get(name) ?? 'not_started';
  }

  public dispose(): void {
    this._onDidChangeTreeData.dispose();
  }
}
