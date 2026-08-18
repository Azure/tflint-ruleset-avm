# terraform_heredoc_usage

Check whether HEREDOC is used for JSON or YAML, if so suggest the user to use the built-in function instead

## Example

```hcl
resource "null_resource" "default" {
  value = <<-JSON
{
    "BigIntSupported": 995815895020119788889,
    "date": "20180322",
    "message": "Success !",
    "status": 200,
    "city": "北京",
    "count": 632,
    "data": {
        "shidu": "34%",
        "pm25": 73,
        "yesterday": {
            "date": "21日星期三",
            "sunrise": "06:19",
            "fl": "<3级",
            "type": "多云",
            "notice": "阴晴之间，谨防紫外线侵扰"
        },
        "forecast": [
            {
                "date": "22日星期四",
                "sunrise": "06:17",
                "notice": "愿你拥有比阳光明媚的心情"
            }
        ]
    }
}
  JSON
}
```

```
$ tflint
1 issue(s) found:

Notice: for JSON, instead of HEREDOC, use a combination of a `local` and the `jsonencode` function (terraform_heredoc_usage)

  on test.tf line 2:
   2:   value = <<-JSON
   3: {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_heredoc_usage.md
```

## Why
Do not use HEREDOC for JSON, YAML since there are better ways to achieve the same outcome using terraform interpolations or resources:
For JSON, use a combination of a local and the jsonencode function
For YAML, use a combination of a local and the yamlencode function
see https://docs.cloudposse.com/reference/best-practices/terraform-best-practices/#do-not-use-heredoc-for-json-yaml-or-iam-policies

## How To Fix
Use the built-in function to parse JSON/YAML instead of HEREDOC