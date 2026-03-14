package model

// EntityType identifies a domain entity kind.
type EntityType string

const (
	EntityCluster           EntityType = "Cluster"
	EntityUser              EntityType = "User"
	EntityRole              EntityType = "Role"
	EntityTrait             EntityType = "Trait"
	EntityNode              EntityType = "Node"
	EntityDatabase          EntityType = "Database"
	EntityDatabaseServer    EntityType = "DatabaseServer"
	EntityKubernetesCluster EntityType = "KubernetesCluster"
	EntityKubernetesServer  EntityType = "KubernetesServer"
	EntityApplication       EntityType = "Application"
	EntityWindowsDesktop    EntityType = "WindowsDesktop"
	EntityTrustedCluster    EntityType = "TrustedCluster"
	EntityLabel             EntityType = "Label"
	EntityLogin             EntityType = "Login"
	EntityNamespace         EntityType = "Namespace"
	EntityPrincipal         EntityType = "Principal"
	EntityResourceSelector  EntityType = "ResourceSelector"
	EntityAccessCapability  EntityType = "AccessCapability"
	EntityRule              EntityType = "Rule"
	EntityCondition         EntityType = "Condition"
)

// RelationshipType identifies a relationship kind in the graph.
type RelationshipType string

const (
	RelUserHasRole          RelationshipType = "USER_HAS_ROLE"
	RelRoleAllowsLogin      RelationshipType = "ROLE_ALLOWS_LOGIN"
	RelRoleMatchesLabel     RelationshipType = "ROLE_MATCHES_LABEL"
	RelUserCanAccessNode    RelationshipType = "USER_CAN_ACCESS_NODE"
	RelUserCanAccessDB      RelationshipType = "USER_CAN_ACCESS_DB"
	RelUserCanAccessKube    RelationshipType = "USER_CAN_ACCESS_KUBE"
	RelUserCanAccessDesktop RelationshipType = "USER_CAN_ACCESS_DESKTOP"
	RelRoleGrantsVerb       RelationshipType = "ROLE_GRANTS_VERB"
	RelResourceInCluster    RelationshipType = "RESOURCE_IN_CLUSTER"
	RelTrustsCluster        RelationshipType = "TRUSTS_CLUSTER"
	RelResourceHasLabel     RelationshipType = "RESOURCE_HAS_LABEL"
	RelUserHasTrait         RelationshipType = "USER_HAS_TRAIT"
	RelRoleRequiresTrait    RelationshipType = "ROLE_REQUIRES_TRAIT"
	RelAccessPath           RelationshipType = "ACCESS_PATH"
	RelPrivilegePath        RelationshipType = "POSSIBLE_PRIVILEGE_PATH"
)
