package ai_analyze_core

import (
	"context"
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/fansqz/fancode-backend/constants"
	"testing"
)

func TestVisualDescriptionAnalyzer_AnalyzeCode(t *testing.T) {
	aiConfig := &config.AIConfig{
		Provider: "volcengine",
		ApiKey:   "e3442936-4b69-4e41-9d28-d1de02326fae",
		ApiBase:  "https://ark.cn-beijing.volces.com/api/v3",
		Model:    "doubao-seed-1.6-250615",
		Timeout:  30,
	}
	visualDescriptionAnalyzer := NewVisualDescriptionAnalyzer(aiConfig)

	// 实现一个二叉树的代码
	cCode := `
	#include <stdio.h>

	struct TreeNode {
		int val;
		struct TreeNode *left;
		struct TreeNode *right;
	};

	int main() {
		struct TreeNode root;
		root.val = 1;
		root.left = NULL;
		root.right = NULL;
		return 0;
	}
	
	`
	// 开启分析
	err := visualDescriptionAnalyzer.StartAnalyzeCode(context.Background(), cCode, constants.LanguageC)
	if err != nil {
		t.Fatalf("AnalyzeCode failed: %v", err)
	}

	// 读取结构体
	visualDescription, err := visualDescriptionAnalyzer.GetVisualDescription(context.Background())
	if err != nil {
		t.Fatalf("GetVisualDescription failed: %v", err)
	}

	t.Logf("VisualDescription: %v", visualDescription)
}
