package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ANIAN0/memos-cli/internal/client"
	"github.com/ANIAN0/memos-cli/pkg/httpclient"
	"github.com/ANIAN0/memos-cli/pkg/output"
)

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "Manage memos",
	Long:  `Create, get, list, update, delete, and search memos on Memos.`,
}

var memoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new memo",
	Long:  `Create a new memo with content and optional visibility/tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		memo, err := mc.Create(cmd.Context(), memoContent, memoVisibility, memoTags)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)
		return out.PrintObject(memo)
	},
}

var memoGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a memo by ID",
	Long:  `Get detailed information about a memo.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		memo, err := mc.Get(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)
		return out.PrintObject(memo)
	},
}

var memoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List memos",
	Long:  `List memos with optional filtering and sorting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		memos, err := mc.List(cmd.Context(), memoPageSize, memoPageToken, memoFilter, memoSort)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)

		items := make([]any, len(memos))
		for i, m := range memos {
			items[i] = m
		}
		return out.PrintList(items)
	},
}

var memoUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a memo",
	Long:  `Update an existing memo's content, visibility, or tags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		memo, err := mc.Update(cmd.Context(), args[0], memoContent, memoVisibility, memoTags)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)
		return out.PrintObject(memo)
	},
}

var memoDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a memo",
	Long:  `Delete a memo by ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		return mc.Delete(cmd.Context(), args[0])
	},
}

var memoSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search memos",
	Long:  `Search for memos containing the query text.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		mc := &client.MemoClient{C: c}
		memos, err := mc.Search(cmd.Context(), args[0], memoPageSize)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)

		items := make([]any, len(memos))
		for i, m := range memos {
			items[i] = m
		}
		return out.PrintList(items)
	},
}

var (
	memoContent    string
	memoVisibility string
	memoTags       []string
	memoPageSize   int
	memoPageToken  string
	memoFilter     string
	memoSort       string
)

func init() {
	memoCreateCmd.Flags().StringVar(&memoContent, "content", "", "memo content (required)")
	memoCreateCmd.Flags().StringVar(&memoVisibility, "visibility", "", "visibility: PRIVATE, PROTECTED, PUBLIC")
	memoCreateCmd.Flags().StringSliceVar(&memoTags, "tags", nil, "tags (comma-separated)")

	memoListCmd.Flags().IntVar(&memoPageSize, "page-size", 0, "page size")
	memoListCmd.Flags().StringVar(&memoPageToken, "page-token", "", "page token")
	memoListCmd.Flags().StringVar(&memoFilter, "filter", "", "filter expression")
	memoListCmd.Flags().StringVar(&memoSort, "sort", "", "sort field")

	memoUpdateCmd.Flags().StringVar(&memoContent, "content", "", "new content")
	memoUpdateCmd.Flags().StringVar(&memoVisibility, "visibility", "", "new visibility")
	memoUpdateCmd.Flags().StringSliceVar(&memoTags, "tags", nil, "new tags")

	memoSearchCmd.Flags().IntVar(&memoPageSize, "page-size", 10, "max results")

	memoCmd.AddCommand(memoCreateCmd)
	memoCmd.AddCommand(memoGetCmd)
	memoCmd.AddCommand(memoListCmd)
	memoCmd.AddCommand(memoUpdateCmd)
	memoCmd.AddCommand(memoDeleteCmd)
	memoCmd.AddCommand(memoSearchCmd)

	rootCmd.AddCommand(memoCmd)
}

// tagsToString converts []string to comma-separated string for display.
func tagsToString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ", ")
}