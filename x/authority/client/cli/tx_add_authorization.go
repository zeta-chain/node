package cli

//func CmdAddAuthorization() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "add-authorization [msg-url] [authorized-policy]",
//		Short: "add a new authorization or update the policy of an existing authorization.Policy type can be 0 for groupEmergency, 1 for groupOperational, 2 for groupAdmin",
//		Args:  cobra.ExactArgs(2),
//		RunE: func(cmd *cobra.Command, args []string) (err error) {
//			clientCtx, err := client.GetClientTxContext(cmd)
//			if err != nil {
//				return err
//			}
//			authorizedPolicy, err := GetPolicyType(args[1])
//
//
//			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
//		},
//	}
//	flags.AddTxFlagsToCmd(cmd)
//	return cmd
//}
//
//func GetPolicyType(policyTypeString string) (types.PolicyType, error) {
//	policyType, err := strconv.ParseInt(policyTypeString, 10, 64)
//	if err != nil {
//		return types.PolicyType_groupEmpty, fmt.Errorf("failed to parse policy type: %w", err)
//	}
//
//	switch policyType {
//	case 0:
//		return types.PolicyType_groupEmergency, nil
//	case 1:
//		return types.PolicyType_groupOperational, nil
//	case 2:
//		return types.PolicyType_groupAdmin, nil
//	default:
//		return types.PolicyType_groupEmpty, fmt.Errorf("invalid policy type: %d", policyType)
//	}
//
//}
