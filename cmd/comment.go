package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ANIAN0/memos-cli/internal/client"
	"github.com/ANIAN0/memos-cli/pkg/httpclient"
	"github.com/ANIAN0/memos-cli/pkg/output"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage memo comments",
	Long:  `List and create comments on memos.`,
}

var commentListCmd = &cobra.Command{
	Use:   "list <memo-id>",
	Short: "List comments on a memo",
	Long:  `List all comments for a specific memo.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		memoID := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		cc := &client.CommentClient{C: c}
		comments, err := cc.List(cmd.Context(), memoID)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)

		items := make([]any, len(comments))
		for i, cm := range comments {
			items[i] = cm
		}
		return out.PrintList(items)
	},
}

var commentCreateCmd = &cobra.Command{
	Use:   "create <memo-id>",
	Short: "Create a comment on a memo",
	Long:  `Create a new comment on a specific memo.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		memoID := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		cc := &client.CommentClient{C: c}
		comment, err := cc.Create(cmd.Context(), memoID, commentContent)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)
		return out.PrintObject(comment)
	},
}

var commentContent string

func init() {
	commentCreateCmd.Flags().StringVar(&commentContent, "content", "", "comment content (required)")
	_ = commentCreateCmd.MarkFlagRequired("content")

	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentCreateCmd)

	rootCmd.AddCommand(commentCmd)
}