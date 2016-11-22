module goroutine_sim();

  reg clk = 1'b0;
  reg reset = 1'b1;
  reg[31:0] counter = 32'h0;
  reg kick = 1'b0;
  wire busy;

  initial begin
    $dumpfile("a.vcd");
    $dumpvars();
  end

  always begin
  #10
  clk <= 1'b0;
  #10
  clk <= 1'b1;
  end

  always @(posedge clk) begin
    counter <= counter + 1;
  end

   wire signed [32-1 : 0] a_out;
   wire signed [32-1 : 0] b_out;
   
   reg signed [32-1 : 0] s1_address;
   reg 			 s1_we;
   reg 			 s1_oe;
   reg signed [32-1 : 0] s1_din;
   reg signed [32-1 : 0] s0_address;
   reg 			 s0_we;
   reg 			 s0_oe;
   reg signed [32-1 : 0] s0_din;
   wire  		 h_busy;
   wire  		 g_busy;
   wire 		 f_busy;

  always @(posedge clk) begin
    if(counter == 5) reset <= 1'b0;
    if(counter == 30) kick <= 1'b1; else kick <= 1'b0;
    if(counter > 32 && busy == 1'b0) begin
       $display(a_out);
       $display(b_out);
       $finish;
    end
  end

  always @(posedge clk) begin
     if(counter >= 6 && counter < 6+16) begin
	s0_oe <= 1'b1;
	s0_we <= 1'b1;
	s0_address <= (counter - 6);
	s0_din <= (counter - 5);
	s1_oe <= 1'b1;
	s1_we <= 1'b1;
	s1_address <= (counter - 6);
	s1_din <= (counter - 5);
     end else begin
	s0_we <= 1'b0;
	s0_oe <= 1'b1;
	s0_address <= 32'h0;
	s0_din <= 32'h0;
	s1_we <= 1'b0;
	s1_oe <= 1'b1;
	s1_address <= 32'h0;
	s1_din <= 32'h0;
     end
  end

   assign busy = h_busy || f_busy || g_busy;
	 
  goroutines U
    (
     .clk(clk),
     .reset(reset),
     .b_in(32'h0), // i
     .b_we(1'b0), // i
     .b_out(b_out), // o
     .a_in(32'h0), // i
     .a_we(1'b0), // i
     .a_out(a_out), // o
     .s1_address(s1_address), // i
     .s1_we(s1_we), // i
     .s1_oe(s1_oe), // i
     .s1_din(s1_din), // i
     .s1_dout(), // o
     .s1_length(), // o
     .s0_address(s0_address), // i
     .s0_we(s0_we), // i
     .s0_oe(s0_oe), // i
     .s0_din(s0_din), // i
     .s0_dout(), // o
     .s0_length(), // o
     .h_busy(h_busy),
     .h_req(kick),
     .g_busy(g_busy),
     .g_req(1'b0),
     .f_busy(f_busy),
     .f_req(1'b0)
     );

endmodule
